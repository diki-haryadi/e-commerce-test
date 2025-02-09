package shopRepository

import (
	"context"
	"fmt"
	shopDomain "github.com/diki-haryadi/go-micro-template/internal/shop/domain"
	shopDto "github.com/diki-haryadi/go-micro-template/internal/shop/dto"
	"github.com/diki-haryadi/ztools/postgres"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepository(conn *postgres.Postgres) shopDomain.Repository {
	return &repository{postgres: conn}
}

func (r *repository) CreateShop(ctx context.Context, req *shopDto.CreateShopRequestDto) (*shopDto.Shop, error) {
	query := `
        INSERT INTO shops (name)
        VALUES ($1)
        RETURNING shop_id, name, created_at, updated_at
    `

	shop := new(shopDto.Shop)
	err := r.postgres.SqlxDB.QueryRowContext(ctx, query, req.Name).
		Scan(&shop.ID, &shop.Name, &shop.CreatedAt, &shop.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating shop: %v", err)
	}

	return shop, nil
}

func (r *repository) GetShopWarehouses(ctx context.Context, shopID int64) ([]*shopDto.Warehouse, error) {
	query := `
        SELECT warehouse_id, name, status, created_at, updated_at
        FROM warehouses
        WHERE shop_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.postgres.SqlxDB.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []*shopDto.Warehouse
	for rows.Next() {
		warehouse := new(shopDto.Warehouse)
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.Name,
			&warehouse.Status,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}
