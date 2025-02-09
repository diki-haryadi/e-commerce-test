package whDto

import (
	"time"
)

type Warehouse struct {
	ID      string `db:"warehouse_id" json:"id"`
	ShopID  string `db:"shop_id" json:"shop_id"`
	Name    string `db:"name" json:"name"`
	Address string `db:"address" json:"address"`
	Status  string `db:"status" json:"status"`
}

type Stock struct {
	ProductID string  `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	Product   Product `db:"product" json:"product"`
}

type WarehouseStock struct {
	ID          string     `db:"id" json:"id"`
	WarehouseID string     `db:"warehouse_id" json:"warehouse_id"`
	ProductID   string     `db:"product_id" json:"product_id"`
	Quantity    int        `db:"quantity" json:"quantity"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	ProductName string     `db:"product_name" json:"product_name"`
	Price       float64    `db:"price" json:"price"`
}

type Product struct {
	Name  string  `db:"name" json:"name"`
	Price float64 `db:"price" json:"price"`
}

type CreateWarehouseRequestDto struct {
	ShopID  string `json:"shop_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type CreateWarehouseResponseDto struct {
	ID      string `json:"id"`
	ShopID  string `json:"shop_id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type UpdateWarehouseStatusRequestDto struct {
	ID          string `json:"id"`
	WarehouseID string `json:"warehouse_id"`
	Status      bool   `json:"status"`
}

type TransferStockRequestDto struct {
	FromWarehouseID string `json:"from_warehouse_id"`
	ToWarehouseID   string `json:"to_warehouse_id"`
	ProductID       string `json:"product_id"`
	Quantity        int    `json:"quantity"`
}

type GetWarehouseStockResponseDto struct {
	WarehouseID string      `json:"warehouse_id"`
	Status      string      `json:"status"`
	Stocks      []StockItem `json:"stocks"`
}

type StockItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
