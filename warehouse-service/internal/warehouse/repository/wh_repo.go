package whRepository

import (
	"context"
	"database/sql"
	"errors"

	whDomain "github.com/diki-haryadi/go-micro-template/internal/warehouse/domain"
	whDto "github.com/diki-haryadi/go-micro-template/internal/warehouse/dto"
	"github.com/diki-haryadi/ztools/postgres"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepository(conn *postgres.Postgres) whDomain.Repository {
	return &repository{postgres: conn}
}

func (r *repository) CreateWarehouse(ctx context.Context, warehouse *whDto.CreateWarehouseRequestDto) (*whDto.Warehouse, error) {
	query := `
        INSERT INTO warehouses (
            shop_id,
            name,
            status
        ) VALUES (
            $1, $2, $3
        )
        RETURNING id, shop_id, name, status, created_at, updated_at
    `

	result := &whDto.Warehouse{}
	err := r.postgres.SqlxDB.QueryRowContext(
		ctx,
		query,
		warehouse.ShopID,
		warehouse.Name,
		true,
	).Scan(
		&result.ID,
		&result.ShopID,
		&result.Name,
		&result.Status,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) UpdateWarehouseStatus(ctx context.Context, warehouseID string, status bool) error {
	query := `
        UPDATE warehouses 
        SET 
            status = $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $2 
        AND deleted_at IS NULL
    `

	result, err := r.postgres.SqlxDB.ExecContext(
		ctx,
		query,
		status,
		warehouseID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("warehouse not found")
	}

	return nil
}

func (r *repository) GetWarehouseByID(ctx context.Context, warehouseID string) (*whDto.Warehouse, error) {
	query := `
        SELECT 
            w.id,
            w.shop_id,
            w.name,
            w.status,
            w.created_at,
            w.updated_at,
            COALESCE(SUM(ws.quantity), 0) as total_stock
        FROM warehouses w
        LEFT JOIN warehouse_stocks ws ON w.id = ws.warehouse_id AND ws.deleted_at IS NULL
        WHERE w.id = $1 
        AND w.deleted_at IS NULL
        GROUP BY w.id, w.shop_id, w.name, w.status, w.created_at, w.updated_at
    `

	warehouse := &whDto.Warehouse{}
	var totalStock int

	err := r.postgres.SqlxDB.QueryRowContext(
		ctx,
		query,
		warehouseID,
	).Scan(
		&warehouse.ID,
		&warehouse.ShopID,
		&warehouse.Name,
		&warehouse.Status,
		&totalStock,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("warehouse not found")
	}

	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (r *repository) TransferStock(ctx context.Context, fromWarehouseID, toWarehouseID, productID string, quantity int) error {
	tx, err := r.postgres.SqlxDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	transferQuery := `
        INSERT INTO stock_transfers (
            from_warehouse_id,
            to_warehouse_id,
            product_id,
            quantity,
            status
        ) VALUES ($1, $2, $3, $4, 'pending')
        RETURNING id
    `
	var transferID string
	err = tx.QueryRowContext(
		ctx,
		transferQuery,
		fromWarehouseID,
		toWarehouseID,
		productID,
		quantity,
	).Scan(&transferID)
	if err != nil {
		return err
	}

	deductQuery := `
        UPDATE warehouse_stocks
        SET 
            quantity = quantity - $1,
            updated_at = CURRENT_TIMESTAMP
        WHERE warehouse_id = $2 
        AND product_id = $3 
        AND quantity >= $1
        AND deleted_at IS NULL
    `
	result, err := tx.ExecContext(ctx, deductQuery, quantity, fromWarehouseID, productID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("insufficient stock in source warehouse")
	}

	addQuery := `
        INSERT INTO warehouse_stocks (
            warehouse_id,
            product_id,
            quantity
        ) VALUES ($1, $2, $3)
        ON CONFLICT (warehouse_id, product_id) WHERE deleted_at IS NULL
        DO UPDATE SET 
            quantity = warehouse_stocks.quantity + $3,
            updated_at = CURRENT_TIMESTAMP
    `
	_, err = tx.ExecContext(ctx, addQuery, toWarehouseID, productID, quantity)
	if err != nil {
		return err
	}

	updateTransferQuery := `
        UPDATE stock_transfers
        SET 
            status = 'completed',
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
    `
	_, err = tx.ExecContext(ctx, updateTransferQuery, transferID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetWarehouseStock(ctx context.Context, warehouseID string) ([]*whDto.WarehouseStock, error) {
	query := `
        SELECT 
            ws.id,
            ws.warehouse_id,
            ws.product_id,
            ws.quantity,
            ws.created_at,
            ws.updated_at,
            p.name as product_name,
            p.price
        FROM warehouse_stocks ws
        JOIN products p ON p.id = ws.product_id AND p.deleted_at IS NULL
        WHERE ws.warehouse_id = $1 
        AND ws.deleted_at IS NULL
    `

	stocks := make([]*whDto.WarehouseStock, 0)
	err := r.postgres.SqlxDB.SelectContext(ctx, &stocks, query, warehouseID)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}
