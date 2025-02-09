package shopUseCase

import (
	"context"
	"fmt"
	sampleExtServiceDomain "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/domain"
	shopDomain "github.com/diki-haryadi/go-micro-template/internal/shop/domain"
	shopDto "github.com/diki-haryadi/go-micro-template/internal/shop/dto"
)

type useCase struct {
	repository              shopDomain.Repository
	kafkaProducer           shopDomain.KafkaProducer
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase
}

func NewUseCase(
	repository shopDomain.Repository,
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase,
	kafkaProducer shopDomain.KafkaProducer,
) shopDomain.UseCase {
	return &useCase{
		repository:              repository,
		kafkaProducer:           kafkaProducer,
		sampleExtServiceUseCase: sampleExtServiceUseCase,
	}
}

func (uc *useCase) CreateShop(ctx context.Context, req *shopDto.CreateShopRequestDto) (*shopDto.CreateShopResponseDto, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("shop name is required")
	}

	shop, err := uc.repository.CreateShop(ctx, req)
	if err != nil {
		return nil, err
	}

	return &shopDto.CreateShopResponseDto{
		ID:   shop.ID,
		Name: shop.Name,
	}, nil
}

func (uc *useCase) GetShopWarehouses(ctx context.Context, shopID int64) (*shopDto.GetShopWarehousesResponseDto, error) {
	if shopID <= 0 {
		return nil, fmt.Errorf("invalid shop ID")
	}

	warehouses, err := uc.repository.GetShopWarehouses(ctx, shopID)
	if err != nil {
		return nil, err
	}

	warehouseItems := make([]shopDto.WarehouseResponseItem, len(warehouses))
	for i, w := range warehouses {
		warehouseItems[i] = shopDto.WarehouseResponseItem{
			ID:     w.ID,
			Name:   w.Name,
			Status: w.Status,
		}
	}

	return &shopDto.GetShopWarehousesResponseDto{
		ShopID:     shopID,
		Warehouses: warehouseItems,
	}, nil
}
