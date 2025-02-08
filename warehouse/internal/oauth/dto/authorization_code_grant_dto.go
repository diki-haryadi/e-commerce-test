package oauthDto

import (
	"fmt"
	oauthModel "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg"
	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

type AuthorizationCodeGrantRequestDto struct {
	GrantType   string `json:"grant_type"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
	ClientID    string `json:"client_id"`
}

func (g *AuthorizationCodeGrantRequestDto) GetFields(ctx echo.Context) *AuthorizationCodeGrantRequestDto {
	return &AuthorizationCodeGrantRequestDto{
		GrantType:   ctx.FormValue("grant_type"),
		Code:        ctx.FormValue("code"),
		RedirectUri: ctx.FormValue("redirect_uri"),
		ClientID:    ctx.FormValue("client_id"),
	}
}

func (g *AuthorizationCodeGrantRequestDto) ToModel(clientID uuid.UUID) *oauthModel.Client {
	return &oauthModel.Client{
		Common: oauthModel.Common{
			ID: clientID,
		},
	}
}

func (caDto *AuthorizationCodeGrantRequestDto) ValidateAuthorizationCodeDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.GrantType,
			validator.Required,
		),
		validator.Field(
			&caDto.Code,
			validator.Required,
		),
		validator.Field(
			&caDto.RedirectUri,
			validator.Required,
		),
	)
}

func NewOauthAuthorizationCode(client *oauthModel.Client, user *oauthModel.Users, expiresIn int, redirectURI, scope string) *oauthModel.AuthorizationCode {
	return &oauthModel.AuthorizationCode{
		Common: oauthModel.Common{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		ClientID:    pkg.StringOrNull(fmt.Sprint(client.ID)),
		UserID:      pkg.StringOrNull(fmt.Sprint(user.ID)),
		Code:        uuid.New().String(),
		ExpiresAt:   time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		RedirectURI: pkg.StringOrNull(redirectURI),
		Scope:       scope,
	}
}
