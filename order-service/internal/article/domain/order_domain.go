package articleDomain

import (
	"context"

	"github.com/google/uuid"

	orderV1 "github.com/diki-haryadi/protobuf-template/go-micro-template/article/v1"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"

	orderDto "github.com/diki-haryadi/go-micro-template/internal/article/dto"
)

type Configurator interface {
	Configure(ctx context.Context) error
}

type UseCase interface {
	Checkout(ctx context.Context, req *orderDto.CheckoutRequestDto) (*orderDto.CheckoutResponseDto, error)
	ReleaseExpiredOrders(ctx context.Context) error
}

type Repository interface {
	CreateOrder(ctx context.Context, order *orderDto.Order, items []*orderDto.OrderItem) (*orderDto.Order, error)
	ReserveStock(ctx context.Context, warehouseID, productID uuid.UUID, quantity int64) error
	ReleaseStock(ctx context.Context, orderID uuid.UUID) error
	GetExpiredOrders(ctx context.Context) ([]*orderDto.Order, error)
	GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]orderDto.Product, error)
}

type GrpcController interface {
	CreateArticle(ctx context.Context, req *orderV1.CreateArticleRequest) (*orderV1.CreateArticleResponse, error)
	GetArticleById(ctx context.Context, req *orderV1.GetArticleByIdRequest) (*orderV1.GetArticleByIdResponse, error)
}

type HttpController interface {
	Checkout(ctx echo.Context) error
}

type Job interface {
	StartJobs(ctx context.Context)
}

type KafkaProducer interface {
	PublishCreateEvent(ctx context.Context, messages ...kafka.Message) error
}

type KafkaConsumer interface {
	RunConsumers(ctx context.Context)
}
