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

type RefreshTokenRequestDto struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (g *RefreshTokenRequestDto) GetFields(ctx echo.Context) *RefreshTokenRequestDto {
	return &RefreshTokenRequestDto{
		GrantType:    ctx.FormValue("grant_type"),
		RefreshToken: ctx.FormValue("refresh_token"),
		Scope:        ctx.FormValue("scope"),
	}
}

func (g *RefreshTokenRequestDto) ToModel(clientID uuid.UUID) *oauthModel.Client {
	return &oauthModel.Client{
		Common: oauthModel.Common{
			ID: clientID,
		},
	}
}

func (caDto *RefreshTokenRequestDto) ValidateRefreshTokenDto() error {
	return validator.ValidateStruct(caDto,
		validator.Field(
			&caDto.GrantType,
			validator.Required,
		),
		validator.Field(
			&caDto.RefreshToken,
			validator.Required,
		),
		validator.Field(
			&caDto.Scope,
			validator.Required,
		),
	)
}

func NewOauthRefreshToken(client *oauthModel.Client, user *oauthModel.Users, expiresIn int, scope string) *oauthModel.RefreshToken {
	refreshToken := &oauthModel.RefreshToken{
		Common: oauthModel.Common{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		ClientID:  pkg.StringOrNull(fmt.Sprint(client.ID)),
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		refreshToken.UserID = pkg.StringOrNull(fmt.Sprint(user.ID))
	}
	return refreshToken
}
