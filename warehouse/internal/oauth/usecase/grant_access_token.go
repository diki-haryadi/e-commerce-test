package oauthUseCase

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg"
)

func (uc *useCase) GrantAccessToken(ctx context.Context, client *oauthDomain.Client, user *oauthDomain.Users, expiresIn int, scope string) (*oauthDto.AccessTokenResponse, error) {
	accessToken, err := uc.repository.GrantAccessToken(ctx, client, user, expiresIn, scope)
	if err != nil {
		return nil, err
	}

	accessTokenResponse, err := oauthDto.NewAccessTokenResponse(
		accessToken,
		nil,
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime, // expires
		pkg.Bearer,
	)
	return accessTokenResponse, nil
}
