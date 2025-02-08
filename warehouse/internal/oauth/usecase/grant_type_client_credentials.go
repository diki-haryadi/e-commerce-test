package oauthUseCase

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
)

func (uc *useCase) ClientCredentialsGrant(ctx context.Context, scope string, refreshToken string, client *oauthDomain.Client) (*oauthDto.AccessTokenResponse, error) {
	// Fetch the refresh token
	theRefreshToken, err := uc.repository.GetValidRefreshToken(refreshToken, client)
	if err != nil {
		return nil, err
	}

	// Get the scope
	scope, err = uc.getRefreshTokenScope(ctx, theRefreshToken, scope)
	if err != nil {
		return nil, err
	}

	// Log in the user
	accessToken, newRefreshToken, err := uc.Login(
		ctx,
		theRefreshToken.Client,
		theRefreshToken.User,
		scope,
	)
	if err != nil {
		return nil, err
	}

	// Create response
	accessTokenResponse, err := oauthDto.NewAccessTokenResponse(
		accessToken,
		newRefreshToken, // refresh token
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime,
		pkg.Bearer,
	)
	if err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}

func (uc *useCase) getRefreshTokenScope(ctx context.Context, refreshToken *oauthDomain.RefreshToken, requestedScope string) (string, error) {
	var (
		scope = refreshToken.Scope // default to the scope originally granted by the resource owner
		err   error
	)

	// If the scope is specified in the request, get the scope string
	if requestedScope != "" {
		scope, err = uc.repository.GetScope(ctx, requestedScope)
		if err != nil {
			return "", err
		}
	}

	// Requested scope CANNOT include any scope not originally granted
	if !pkg.SpaceDelimitedStringNotGreater(scope, refreshToken.Scope) {
		return "", response.ErrRequestedScopeCannotBeGreater
	}

	return scope, nil
}
