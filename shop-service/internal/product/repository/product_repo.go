package productRepository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"math"

	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
	productDto "github.com/diki-haryadi/go-micro-template/internal/product/dto"
	"github.com/diki-haryadi/ztools/postgres"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepository(conn *postgres.Postgres) productDomain.Repository {
	return &repository{postgres: conn}
}

func (rp *repository) GetProducts(ctx context.Context, filter *productDto.GetProductsRequestDto) (*productDto.GetProductsResponseDto, error) {
	var args []interface{}
	argCount := 1

	query := `
        WITH filtered_products AS (
            SELECT 
                id, name, description, category_id, 
                price, status
            FROM products
            WHERE 1=1
    `

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filter.Search+"%")
		argCount++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filter.Status)
		argCount++
	}

	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, *filter.MinPrice)
		argCount++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, *filter.MaxPrice)
		argCount++
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM filtered_products"
	var totalRecords int64
	err := rp.postgres.SqlxDB.GetContext(ctx, &totalRecords, query+") "+countQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error counting products: %w", err)
	}

	// Add pagination
	offset := (filter.Page - 1) * filter.PageSize
	query += fmt.Sprintf(`) SELECT * FROM filtered_products 
        ORDER BY name 
        LIMIT $%d OFFSET $%d`, argCount, argCount+1)
	args = append(args, filter.PageSize, offset)

	// Execute final query
	rows, err := rp.postgres.SqlxDB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []productDto.Products
	//productDto.ProductResponse
	for rows.Next() {
		var product productDto.Products
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.CategoryID,
			&product.Price,
			&product.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(filter.PageSize)))

	return &productDto.GetProductsResponseDto{
		Products: products,
		Meta: productDto.PaginationMeta{
			CurrentPage:  filter.Page,
			PageSize:     filter.PageSize,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	}, nil
}

func (rp *repository) GetProductByID(ctx context.Context, id uuid.UUID) (*productDto.ProductResponse, error) {
	query := `
        SELECT 
            id, name, description, category_id, price, status
        FROM products 
        WHERE id = $1`

	var product productDto.Products
	err := rp.postgres.SqlxDB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.Price,
		&product.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("error querying product: %w", err)
	}
	resp := productDto.ProductResponse(product)
	return &resp, nil
}

func (r *repository) GetProductStock(ctx context.Context, productIDs []string) (map[string]int64, error) {
	query := `
        SELECT ws.product_id, COALESCE(SUM(ws.quantity), 0) as total_stock
        FROM warehouse_stocks ws
        JOIN warehouses w ON w.warehouse_id = ws.warehouse_id
        WHERE ws.product_id = ANY($1) AND w.status = 'active'
        GROUP BY ws.product_id
    `

	rows, err := r.postgres.SqlxDB.QueryContext(ctx, query, pq.Array(productIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stocks := make(map[string]int64)
	for rows.Next() {
		var productID string
		var stock int64
		if err := rows.Scan(&productID, &stock); err != nil {
			return nil, err
		}
		stocks[productID] = stock
	}

	return stocks, nil
}
