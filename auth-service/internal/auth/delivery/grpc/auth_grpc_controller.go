package articleGrpcController

import (
	articleDomain "github.com/diki-haryadi/go-micro-template/internal/auth/domain"
)

type controller struct {
	useCase articleDomain.UseCase
}

func NewController(uc articleDomain.UseCase) articleDomain.GrpcController {
	return &controller{
		useCase: uc,
	}
}

//func (c *controller) CreateArticle(ctx context.Context, req *articleV1.CreateArticleRequest) (*articleV1.CreateArticleResponse, error) {
//	aDto := &articleDto.SignUpRequestDto{
//		Username: req.Name,
//		Password: req.Desc,
//	}
//	err := aDto.ValidateCreateArticleDto()
//	if err != nil {
//		return nil, articleException.CreateArticleValidationExc(err)
//	}
//
//	article, err := c.useCase.SignUp(ctx, aDto)
//	if err != nil {
//		return nil, err
//	}
//
//	return &articleV1.CreateArticleResponse{
//		Id:   article.ID.String(),
//		Name: article.Username,
//		Desc: article.Password,
//	}, nil
//}
//
//func (c *controller) GetArticleById(ctx context.Context, req *articleV1.GetArticleByIdRequest) (*articleV1.GetArticleByIdResponse, error) {
//	return nil, status.Error(codes.Unimplemented, "not implemented")
//}
