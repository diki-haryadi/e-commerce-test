package productDto

import (
	"github.com/google/uuid"
)

type Products struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CategoryID  uuid.UUID `db:"category_id" json:"category_id"`
	Price       int64     `db:"price" json:"price"`
	Stock       int64     `db:"stock" json:"stock"`
	Status      bool      `db:"status" json:"status"`
}

type ProductStock struct {
	ProductID   int64  `db:"product_id"`
	Stock       int    `db:"stock"`
	ProductName string `db:"name"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"category_id"`
	Price       int64     `json:"price"`
	Stock       int64     `json:"stock"`
	Status      bool      `json:"status"`
}

type GetProductsRequestDto struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Search   string `query:"search"`
	Status   *bool  `query:"status"`
	MinPrice *int64 `query:"min_price"`
	MaxPrice *int64 `query:"max_price"`
}

type GetProductsResponseDto struct {
	Products []Products     `json:"products"`
	Meta     PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
}

type GetProductByIDRequestDto struct {
	ID uuid.UUID `param:"id"`
}
