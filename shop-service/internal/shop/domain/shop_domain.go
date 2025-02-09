package shopDomain

import (
	"context"

	shopDto "github.com/diki-haryadi/go-micro-template/internal/shop/dto"
	"github.com/labstack/echo/v4"
)

type Configurator interface {
	Configure(ctx context.Context) error
}

type UseCase interface {
	CreateShop(ctx context.Context, req *shopDto.CreateShopRequestDto) (*shopDto.CreateShopResponseDto, error)
	GetShopWarehouses(ctx context.Context, shopID int64) (*shopDto.GetShopWarehousesResponseDto, error)
}

type Repository interface {
	CreateShop(ctx context.Context, shop *shopDto.CreateShopRequestDto) (*shopDto.Shop, error)
	GetShopWarehouses(ctx context.Context, shopID int64) ([]*shopDto.Warehouse, error)
}

type GrpcController interface {
	//CreateArticle(ctx context.Context, req *productV1.CreateArticleRequest) (*productV1.CreateArticleResponse, error)
	//GetArticleById(ctx context.Context, req *productV1.GetArticleByIdRequest) (*productV1.GetArticleByIdResponse, error)
}

type HttpController interface {
	CreateShop(c echo.Context) error
	GetShopWarehouses(c echo.Context) error
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
