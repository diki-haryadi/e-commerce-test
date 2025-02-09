package articleRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"time"

	orderDomain "github.com/diki-haryadi/go-micro-template/internal/order/domain"
	orderDto "github.com/diki-haryadi/go-micro-template/internal/order/dto"
	"github.com/diki-haryadi/ztools/postgres"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepository(conn *postgres.Postgres) orderDomain.Repository {
	return &repository{postgres: conn}
}

func (r *repository) CreateOrder(ctx context.Context, order *orderDto.Order, items []*orderDto.OrderItem) (*orderDto.Order, error) {
	tx, err := r.postgres.SqlxDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert order
	query := `
        INSERT INTO orders (id, user_id, status, total_amount, payment_deadline, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

	err = tx.QueryRowContext(ctx, query,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalAmount,
		order.PaymentDeadline,
	).Scan(&order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Insert order items and update stock directly
	for _, item := range items {
		// First check and update stock
		stockQuery := `
            UPDATE warehouse_stocks 
            SET quantity = quantity - $1, 
                updated_at = NOW()
            WHERE warehouse_id = $2 
            AND product_id = $3 
            AND quantity >= $1
            RETURNING quantity
        `

		var remainingStock int64
		err = tx.QueryRowContext(ctx, stockQuery,
			item.Quantity,
			item.WarehouseID,
			item.ProductID,
		).Scan(&remainingStock)

		if err == sql.ErrNoRows {
			return nil, errors.New("insufficient stock")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to update stock: %v", err)
		}

		// Then insert order item
		itemQuery := `
            INSERT INTO order_items (
                id, order_id, product_id, warehouse_id, 
                quantity, price, stock_status, created_at, updated_at
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `

		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID,
			order.ID,
			item.ProductID,
			item.WarehouseID,
			item.Quantity,
			item.Price,
			item.StockStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return order, nil
}

func (r *repository) ReleaseStock(ctx context.Context, orderID uuid.UUID) error {
	tx, err := r.postgres.SqlxDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update order status
	orderQuery := `
        UPDATE orders 
        SET status = 'cancelled', 
            updated_at = NOW()
        WHERE id = $1 AND status = 'pending'
    `

	result, err := tx.ExecContext(ctx, orderQuery, orderID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("order not found or already processed")
	}

	// Return stock to warehouse
	stockQuery := `
        UPDATE warehouse_stocks ws
        SET quantity = ws.quantity + oi.quantity,
            updated_at = NOW()
        FROM order_items oi
        WHERE oi.order_id = $1
        AND ws.warehouse_id = oi.warehouse_id
        AND ws.product_id = oi.product_id
        AND oi.stock_status = 'reserved'
    `

	_, err = tx.ExecContext(ctx, stockQuery, orderID)
	if err != nil {
		return err
	}

	// Update order items status
	itemQuery := `
        UPDATE order_items
        SET stock_status = 'released',
            updated_at = NOW()
        WHERE order_id = $1
        AND stock_status = 'reserved'
    `

	_, err = tx.ExecContext(ctx, itemQuery, orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetExpiredOrders(ctx context.Context) ([]*orderDto.Order, error) {
	query := `
        SELECT id, user_id, status, total_amount, payment_deadline, created_at, updated_at
        FROM orders
        WHERE status = 'pending'
        AND payment_deadline < NOW()
    `

	rows, err := r.postgres.SqlxDB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*orderDto.Order
	for rows.Next() {
		order := new(orderDto.Order)
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.PaymentDeadline,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *repository) ReserveStock(ctx context.Context, warehouseID, productID uuid.UUID, quantity int64) error {
	// Start transaction
	tx, err := r.postgres.SqlxDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the stock record for update
	query := `
       SELECT quantity 
       FROM warehouse_stocks 
       WHERE warehouse_id = $1 
       AND product_id = $2
       FOR UPDATE
   `

	var currentStock int64
	err = tx.QueryRowContext(ctx, query, warehouseID, productID).Scan(&currentStock)
	if err == sql.ErrNoRows {
		return errors.New("stock not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get stock: %v", err)
	}

	// Check if stock is sufficient
	if currentStock < quantity {
		return errors.New("insufficient stock")
	}

	// Update stock
	updateQuery := `
       UPDATE warehouse_stocks 
       SET quantity = quantity - $1,
           updated_at = NOW()
       WHERE warehouse_id = $2 
       AND product_id = $3
   `

	result, err := tx.ExecContext(ctx, updateQuery, quantity, warehouseID, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("stock not found")
	}

	// Add stock reservation record
	reserveQuery := `
       INSERT INTO stock_reservations (
           id,
           warehouse_id,
           product_id,
           quantity,
           status,
           created_at,
           updated_at
       ) VALUES ($1, $2, $3, $4, $5, $6, $6)
   `

	_, err = tx.ExecContext(ctx, reserveQuery,
		uuid.New(),
		warehouseID,
		productID,
		quantity,
		"reserved",
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create reservation record: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *repository) GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]orderDto.Product, error) {
	query := `
        SELECT id, name, price
        FROM products
        WHERE id = ANY($1)
    `

	rows, err := r.postgres.SqlxDB.QueryContext(ctx, query, pq.Array(productIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	products := make(map[uuid.UUID]orderDto.Product)
	for rows.Next() {
		var p orderDto.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products[p.ID] = p
	}

	return products, nil
}
