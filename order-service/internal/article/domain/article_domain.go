package articleDomain

import (
	"context"

	"github.com/google/uuid"

	articleV1 "github.com/diki-haryadi/protobuf-template/go-micro-template/article/v1"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"

	articleDto "github.com/diki-haryadi/go-micro-template/internal/article/dto"
)

type Article struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"desc"`
}

type Configurator interface {
	Configure(ctx context.Context) error
}

type UseCase interface {
	Checkout(ctx context.Context, req *articleDto.CheckoutRequestDto) (*articleDto.CheckoutResponseDto, error)
	ReleaseExpiredOrders(ctx context.Context) error
}

type Repository interface {
	CreateOrder(ctx context.Context, order *articleDto.Order, items []*articleDto.OrderItem) (*articleDto.Order, error)
	ReserveStock(ctx context.Context, warehouseID, productID uuid.UUID, quantity int64) error
	ReleaseStock(ctx context.Context, orderID uuid.UUID) error
	GetExpiredOrders(ctx context.Context) ([]*articleDto.Order, error)
	GetProductsByIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]articleDto.Product, error)
}

type GrpcController interface {
	CreateArticle(ctx context.Context, req *articleV1.CreateArticleRequest) (*articleV1.CreateArticleResponse, error)
	GetArticleById(ctx context.Context, req *articleV1.GetArticleByIdRequest) (*articleV1.GetArticleByIdResponse, error)
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
