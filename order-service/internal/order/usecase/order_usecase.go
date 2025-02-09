package articleUseCase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"

	sampleExtServiceDomain "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/domain"
	orderDomain "github.com/diki-haryadi/go-micro-template/internal/order/domain"
	orderDto "github.com/diki-haryadi/go-micro-template/internal/order/dto"
)

type useCase struct {
	repository              orderDomain.Repository
	kafkaProducer           orderDomain.KafkaProducer
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase
}

func NewUseCase(
	repository orderDomain.Repository,
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase,
	kafkaProducer orderDomain.KafkaProducer,
) orderDomain.UseCase {
	return &useCase{
		repository:              repository,
		kafkaProducer:           kafkaProducer,
		sampleExtServiceUseCase: sampleExtServiceUseCase,
	}
}

func (uc *useCase) Checkout(ctx context.Context, req *orderDto.CheckoutRequestDto) (*orderDto.CheckoutResponseDto, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	products, err := uc.repository.GetProductsByIDs(ctx, extractProductIDs(req.Items))
	if err != nil {
		return nil, err
	}

	if len(products) != len(req.Items) {
		return nil, errors.New("some products not found")
	}

	totalAmount := calculateTotalAmount(products, req.Items)

	order := &orderDto.Order{
		ID:              uuid.New(),
		UserID:          req.UserID,
		Status:          "pending",
		TotalAmount:     totalAmount,
		PaymentDeadline: time.Now().Add(1 * time.Hour),
	}

	// Create order items
	orderItems := make([]*orderDto.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderItems[i] = &orderDto.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			WarehouseID: item.WarehouseID,
			Quantity:    item.Quantity,
			Price:       products[item.ProductID].Price,
			StockStatus: "reserved",
		}
	}

	createdOrder, err := uc.repository.CreateOrder(ctx, order, orderItems)
	if err != nil {
		return nil, err
	}

	return &orderDto.CheckoutResponseDto{
		OrderID:         createdOrder.ID,
		Status:          createdOrder.Status,
		TotalAmount:     createdOrder.TotalAmount,
		PaymentDeadline: createdOrder.PaymentDeadline,
	}, nil
}

func (uc *useCase) ReleaseExpiredOrders(ctx context.Context) error {
	orders, err := uc.repository.GetExpiredOrders(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := uc.repository.ReleaseStock(ctx, order.ID); err != nil {
			log.Printf("Failed to release stock for order %s: %v", order.ID, err)
			continue
		}
	}

	return nil
}

func extractProductIDs(items []orderDto.OrderItemRequest) []uuid.UUID {
	productIDs := make([]uuid.UUID, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID
	}
	return productIDs
}

func calculateTotalAmount(products map[uuid.UUID]orderDto.Product, items []orderDto.OrderItemRequest) int64 {
	var total int64
	for _, item := range items {
		if product, ok := products[item.ProductID]; ok {
			total += product.Price * item.Quantity
		}
	}
	return total
}
