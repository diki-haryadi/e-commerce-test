package productDomain

import (
	"context"

	"github.com/google/uuid"

	productDto "github.com/diki-haryadi/go-micro-template/internal/product/dto"
	"github.com/labstack/echo/v4"
)

type Configurator interface {
	Configure(ctx context.Context) error
}

type UseCase interface {
	GetProducts(ctx context.Context, req *productDto.GetProductsRequestDto) (*productDto.GetProductsResponseDto, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*productDto.ProductResponse, error)
}

type Repository interface {
	GetProducts(ctx context.Context, filter *productDto.GetProductsRequestDto) (*productDto.GetProductsResponseDto, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*productDto.ProductResponse, error)
	GetProductStock(ctx context.Context, productIDs []string) (map[string]int64, error)
}

type GrpcController interface {
	//CreateArticle(ctx context.Context, req *productV1.CreateArticleRequest) (*productV1.CreateArticleResponse, error)
	//GetArticleById(ctx context.Context, req *productV1.GetArticleByIdRequest) (*productV1.GetArticleByIdResponse, error)
}

type HttpController interface {
	GetProducts(ctx echo.Context) error
	GetProductByID(ctx echo.Context) error
}

type Job interface {
	//StartJobs(ctx context.Context)
}

type KafkaProducer interface {
	//PublishCreateEvent(ctx context.Context, messages ...kafka.Message) error
}

type KafkaConsumer interface {
	//RunConsumers(ctx context.Context)
}
