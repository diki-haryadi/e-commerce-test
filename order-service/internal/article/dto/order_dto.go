package articleDto

import (
	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID              uuid.UUID `db:"id"`
	UserID          uuid.UUID `db:"user_id"`
	Status          string    `db:"status"`
	TotalAmount     int64     `db:"total_amount"`
	PaymentDeadline time.Time `db:"payment_deadline"`
}

type OrderItem struct {
	ID          uuid.UUID `db:"id"`
	OrderID     uuid.UUID `db:"order_id"`
	ProductID   uuid.UUID `db:"product_id"`
	WarehouseID uuid.UUID `db:"warehouse_id"`
	Quantity    int64     `db:"quantity"`
	Price       int64     `db:"price"`
	StockStatus string    `db:"stock_status"`
}

type OrderItemRequest struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	WarehouseID uuid.UUID `json:"warehouse_id" validate:"required"`
	Quantity    int64     `json:"quantity" validate:"required,gt=0"`
}

type CheckoutRequestDto struct {
	UserID uuid.UUID          `json:"user_id" validate:"required"`
	Items  []OrderItemRequest `json:"items" validate:"required,dive"`
}

type CheckoutResponseDto struct {
	OrderID         uuid.UUID `json:"order_id"`
	Status          string    `json:"status"`
	TotalAmount     int64     `json:"total_amount"`
	PaymentDeadline time.Time `json:"payment_deadline"`
}

func (r *CheckoutRequestDto) Validate() error {
	return validator.ValidateStruct(r,
		validator.Field(&r.UserID, validator.Required),
		validator.Field(&r.Items, validator.Required, validator.Length(1, 100)),
	)
}

type Product struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Price int64     `db:"price"`
}
