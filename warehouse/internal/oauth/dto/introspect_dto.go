package oauthDto

import (
	oauthModel "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type IntrospectRequestDto struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
}

func (g *IntrospectRequestDto) GetFields(ctx echo.Context) *IntrospectRequestDto {
	return &IntrospectRequestDto{
		Token:         ctx.FormValue("token"),
		TokenTypeHint: ctx.FormValue("token_type_hint"),
	}
}

func (g *IntrospectRequestDto) ToModel(clientID uuid.UUID) *oauthModel.Client {
	return &oauthModel.Client{
		Common: oauthModel.Common{
			ID: clientID,
		},
	}
}

func (caDto *IntrospectRequestDto) ValidateIntrospectDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.Token,
			validator.Required,
		),
		validator.Field(
			&caDto.TokenTypeHint,
			validator.Required,
		),
	)
}

// IntrospectResponse ...
type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Username  string `json:"username,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	ExpiresAt int    `json:"exp,omitempty"`
}
