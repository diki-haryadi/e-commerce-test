package oauthDto

import (
	oauthModel "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ClientCredentialsGrantRequestDto struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`
}

func (g *ClientCredentialsGrantRequestDto) GetFields(ctx echo.Context) *ClientCredentialsGrantRequestDto {
	return &ClientCredentialsGrantRequestDto{
		GrantType: ctx.FormValue("grant_type"),
		Scope:     ctx.FormValue("scope"),
	}
}

func (g *ClientCredentialsGrantRequestDto) ToModel(clientID uuid.UUID) *oauthModel.Client {
	return &oauthModel.Client{
		Common: oauthModel.Common{
			ID: clientID,
		},
	}
}

func (caDto *ClientCredentialsGrantRequestDto) ValidateClientCredentialsDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.GrantType,
			validator.Required,
		),
		validator.Field(
			&caDto.Scope,
			validator.Required,
		),
	)
}
