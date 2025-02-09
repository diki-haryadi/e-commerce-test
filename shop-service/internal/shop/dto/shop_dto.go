package shopDto

import (
	validator "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Shop struct {
	ID        int64     `db:"shop_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateShopRequestDto struct {
	Name string `json:"name"`
}

func (r *CreateShopRequestDto) Validate() error {
	return validator.ValidateStruct(r,
		validator.Field(&r.Name, validator.Required, validator.Length(3, 255)),
	)
}

type CreateShopResponseDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetShopWarehousesResponseDto struct {
	ShopID     int64                   `json:"shop_id"`
	Warehouses []WarehouseResponseItem `json:"warehouses"`
}

type WarehouseResponseItem struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Warehouse struct {
	ID        int64     `db:"warehouse_id" json:"id"`
	ShopID    int64     `db:"shop_id" json:"shop_id"`
	Name      string    `db:"name" json:"name"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
