package oauthGrpcController

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	articleV1 "github.com/diki-haryadi/protobuf-template/go-micro-template/article/v1"

	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
)

type controller struct {
	useCase oauthDomain.UseCase
}

func NewController(uc oauthDomain.UseCase) oauthDomain.GrpcController {
	return &controller{
		useCase: uc,
	}
}

func (c *controller) CreateArticle(ctx context.Context, req *articleV1.CreateArticleRequest) (*articleV1.CreateArticleResponse, error) {
	//aDto := &oauthDto.{
	//	Name:        req.Name,
	//	Description: req.Desc,
	//}
	//err := aDto.ValidateCreateArticleDto()
	//if err != nil {
	//	return nil, oauthException.CreateArticleValidationExc(err)
	//}
	//
	//article, err := c.useCase.CreateArticle(ctx, aDto)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &articleV1.CreateArticleResponse{
	//	Id:   article.ID.String(),
	//	Name: article.Name,
	//	Desc: article.Description,
	//}, nil
	return nil, nil
}

func (c *controller) GetArticleById(ctx context.Context, req *articleV1.GetArticleByIdRequest) (*articleV1.GetArticleByIdResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
