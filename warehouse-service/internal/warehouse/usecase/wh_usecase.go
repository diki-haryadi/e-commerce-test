package whUseCase

import (
	"context"
	"errors"
	sampleExtServiceDomain "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/domain"
	whDomain "github.com/diki-haryadi/go-micro-template/internal/warehouse/domain"
	whDto "github.com/diki-haryadi/go-micro-template/internal/warehouse/dto"
)

type useCase struct {
	repository              whDomain.Repository
	kafkaProducer           whDomain.KafkaProducer
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase
}

func NewUseCase(
	repository whDomain.Repository,
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase,
	kafkaProducer whDomain.KafkaProducer,
) whDomain.UseCase {
	return &useCase{
		repository:              repository,
		kafkaProducer:           kafkaProducer,
		sampleExtServiceUseCase: sampleExtServiceUseCase,
	}
}
func (u *useCase) CreateWarehouse(ctx context.Context, req *whDto.CreateWarehouseRequestDto) (*whDto.CreateWarehouseResponseDto, error) {
	warehouse, err := u.repository.CreateWarehouse(ctx, req)
	if err != nil {
		return nil, err
	}

	return &whDto.CreateWarehouseResponseDto{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Address: warehouse.Address,
		Status:  warehouse.Status,
	}, nil
}

func (u *useCase) UpdateWarehouseStatus(ctx context.Context, req *whDto.UpdateWarehouseStatusRequestDto) error {
	return u.repository.UpdateWarehouseStatus(ctx, req.ID, req.Status)
}

func (u *useCase) TransferStock(ctx context.Context, req *whDto.TransferStockRequestDto) error {
	sourceWarehouse, err := u.repository.GetWarehouseByID(ctx, req.FromWarehouseID)
	if err != nil {
		return err
	}
	if sourceWarehouse.Status != "active" {
		return errors.New("source warehouse is not active")
	}

	destWarehouse, err := u.repository.GetWarehouseByID(ctx, req.ToWarehouseID)
	if err != nil {
		return err
	}
	if destWarehouse.Status != "active" {
		return errors.New("destination warehouse is not active")
	}

	return u.repository.TransferStock(ctx, req.FromWarehouseID, req.ToWarehouseID, req.ProductID, req.Quantity)
}

func (u *useCase) GetWarehouseStock(ctx context.Context, warehouseID string) (*whDto.GetWarehouseStockResponseDto, error) {
	warehouse, err := u.repository.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	stocks, err := u.repository.GetWarehouseStock(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	stockItems := make([]whDto.StockItem, 0, len(stocks))
	for _, stock := range stocks {
		stockItems = append(stockItems, whDto.StockItem{
			ProductID: stock.ProductID,
			Name:      stock.ProductName,
			Quantity:  stock.Quantity,
			Price:     stock.Price,
		})
	}

	return &whDto.GetWarehouseStockResponseDto{
		WarehouseID: warehouseID,
		Status:      warehouse.Status,
		Stocks:      stockItems,
	}, nil
}
