package oauthUseCase

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
)

// Login creates an access token and refresh token for a user (logs him/her in)
func (uc *useCase) Login(ctx context.Context, client *oauthDomain.Client, user *oauthDomain.Users, scope string) (*oauthDomain.AccessToken, *oauthDomain.RefreshToken, error) {
	if !uc.IsRoleAllowed(user.Role.Name) {
		return nil, nil, response.ErrInvalidUsernameOrPassword
	}

	// Create a new access token
	accessToken, err := uc.repository.GrantAccessToken(
		ctx,
		client,
		user,
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime, // expires in
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	// Create or retrieve a refresh token
	refreshToken, err := uc.repository.GetOrCreateRefreshToken(
		client,
		user,
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime, // expires in
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
