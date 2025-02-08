package oauthUseCase

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg"
)

func (uc *useCase) RefreshTokenGrant(ctx context.Context, token, scope string, client *oauthDomain.Client) (*oauthDto.AccessTokenResponse, error) {
	// Fetch the refresh token
	theRefreshToken, err := uc.repository.GetValidRefreshToken(token, client)
	if err != nil {
		return nil, err
	}

	// Get the scope
	scopeR, err := uc.GetRefreshTokenScope(ctx, theRefreshToken, scope)
	if err != nil {
		return nil, err
	}

	// Log in the user
	accessToken, refreshToken, err := uc.Login(ctx,
		theRefreshToken.Client,
		theRefreshToken.User,
		scopeR,
	)
	if err != nil {
		return nil, err
	}

	// Create response
	accessTokenResponse, err := oauthDto.NewAccessTokenResponse(
		accessToken,
		refreshToken,
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime, // expires
		pkg.Bearer,
	)
	if err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}
