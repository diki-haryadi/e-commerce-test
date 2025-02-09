package productGrpcController

import (
	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
)

type controller struct {
	useCase productDomain.UseCase
}

func NewController(uc productDomain.UseCase) productDomain.GrpcController {
	return &controller{
		useCase: uc,
	}
}

//func (c *controller) CreateArticle(ctx context.Context, req *articleV1.CreateArticleRequest) (*articleV1.CreateArticleResponse, error) {
//	aDto := &articleDto.CreateArticleRequestDto{
//		Name:        req.Name,
//		Description: req.Desc,
//	}
//	err := aDto.ValidateCreateArticleDto()
//	if err != nil {
//		return nil, articleException.CreateArticleValidationExc(err)
//	}
//
//	article, err := c.useCase.CreateArticle(ctx, aDto)
//	if err != nil {
//		return nil, err
//	}
//
//	return &articleV1.CreateArticleResponse{
//		ID:   article.ID.String(),
//		Name: article.Name,
//		Desc: article.Description,
//	}, nil
//}
//
//func (c *controller) GetArticleById(ctx context.Context, req *articleV1.GetArticleByIdRequest) (*articleV1.GetArticleByIdResponse, error) {
//	return nil, status.Error(codes.Unimplemented, "not implemented")
//}
