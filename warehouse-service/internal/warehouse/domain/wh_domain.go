package productDomain

import (
	"context"

	warehouseDto "github.com/diki-haryadi/go-micro-template/internal/warehouse/dto"
	"github.com/labstack/echo/v4"
)

type Configurator interface {
	Configure(ctx context.Context) error
}

type UseCase interface {
	CreateWarehouse(ctx context.Context, req *warehouseDto.CreateWarehouseRequestDto) (*warehouseDto.CreateWarehouseResponseDto, error)
	UpdateWarehouseStatus(ctx context.Context, req *warehouseDto.UpdateWarehouseStatusRequestDto) error
	TransferStock(ctx context.Context, req *warehouseDto.TransferStockRequestDto) error
	GetWarehouseStock(ctx context.Context, warehouseID string) (*warehouseDto.GetWarehouseStockResponseDto, error)
}

type Repository interface {
	CreateWarehouse(ctx context.Context, warehouse *warehouseDto.CreateWarehouseRequestDto) (*warehouseDto.Warehouse, error)
	UpdateWarehouseStatus(ctx context.Context, warehouseID string, status bool) error
	GetWarehouseStock(ctx context.Context, warehouseID string) ([]*warehouseDto.WarehouseStock, error)
	TransferStock(ctx context.Context, fromWarehouseID, toWarehouseID, productID string, quantity int) error
	GetWarehouseByID(ctx context.Context, warehouseID string) (*warehouseDto.Warehouse, error)
}

type GrpcController interface {
	//CreateArticle(ctx context.Context, req *productV1.CreateArticleRequest) (*productV1.CreateArticleResponse, error)
	//GetArticleById(ctx context.Context, req *productV1.GetArticleByIdRequest) (*productV1.GetArticleByIdResponse, error)
}

type HttpController interface {
	CreateWarehouse(ctx echo.Context) error
	UpdateWarehouseStatus(ctx echo.Context) error
	TransferStock(ctx echo.Context) error
	GetWarehouseStock(ctx echo.Context) error
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
