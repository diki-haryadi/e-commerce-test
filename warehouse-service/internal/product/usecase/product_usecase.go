package productUseCase

import (
	"context"
	"fmt"
	sampleExtServiceDomain "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/domain"
	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
	productDto "github.com/diki-haryadi/go-micro-template/internal/product/dto"
	"github.com/google/uuid"
)

type useCase struct {
	repository              productDomain.Repository
	kafkaProducer           productDomain.KafkaProducer
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase
}

func NewUseCase(
	repository productDomain.Repository,
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase,
	kafkaProducer productDomain.KafkaProducer,
) productDomain.UseCase {
	return &useCase{
		repository:              repository,
		kafkaProducer:           kafkaProducer,
		sampleExtServiceUseCase: sampleExtServiceUseCase,
	}
}

func (uc *useCase) GetProducts(ctx context.Context, req *productDto.GetProductsRequestDto) (*productDto.GetProductsResponseDto, error) {
	// Set default values for pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// Add maximum page size to prevent large queries
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	resp, err := uc.repository.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	// Get product IDs for stock lookup
	productIDs := make([]string, len(resp.Products))
	for i, p := range resp.Products {
		productIDs[i] = p.ID.String()
	}

	// Get stock information
	stocks, err := uc.repository.GetProductStock(ctx, productIDs)
	if err != nil {
		return nil, err
	}

	// Build response
	items := make([]productDto.Products, len(resp.Products))
	for i, p := range resp.Products {
		items[i] = productDto.Products{
			ID:          p.ID,
			Name:        p.Name,
			CategoryID:  p.CategoryID,
			Description: p.Description,
			Price:       p.Price,
			Stock:       stocks[p.ID.String()],
		}
	}

	return &productDto.GetProductsResponseDto{
		Products: items,
		Meta:     resp.Meta,
	}, nil
}

func (uc *useCase) GetProductByID(ctx context.Context, id uuid.UUID) (*productDto.ProductResponse, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid product ID")
	}

	return uc.repository.GetProductByID(ctx, id)
}
