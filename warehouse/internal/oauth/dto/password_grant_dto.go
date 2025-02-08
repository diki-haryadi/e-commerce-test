package oauthDto

import (
	oauthModel "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PasswordGrantRequestDto struct {
	GrantType string `json:"grant_type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Scope     string `json:"scope"`
}

func (g *PasswordGrantRequestDto) GetFields(ctx echo.Context) *PasswordGrantRequestDto {
	return &PasswordGrantRequestDto{
		GrantType: ctx.FormValue("grant_type"),
		Username:  ctx.FormValue("username"),
		Password:  ctx.FormValue("password"),
		Scope:     ctx.FormValue("scope"),
	}
}

func (g *PasswordGrantRequestDto) ToModel(clientID uuid.UUID) *oauthModel.Client {
	return &oauthModel.Client{
		Common: oauthModel.Common{
			ID: clientID,
		},
	}
}

func (caDto *PasswordGrantRequestDto) ValidatePasswordDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.GrantType,
			validator.Required,
		),
		validator.Field(
			&caDto.Username,
			validator.Required,
		),
		validator.Field(
			&caDto.Password,
			validator.Required,
		),
		validator.Field(
			&caDto.Scope,
			validator.Required,
		),
	)
}
