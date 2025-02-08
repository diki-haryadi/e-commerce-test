package oauthUseCase

import (
	"context"
	"fmt"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"time"
)

func (uc *useCase) AuthorizationCodeGrant(ctx context.Context, code, redirectURI string, client *oauthDomain.Client) (*oauthDto.AccessTokenResponse, error) {
	// 1. Fetch the authorization code from the database
	authorizationCode, err := uc.repository.FetchAuthorizationCodeByCode(ctx, client, code)
	if err != nil {
		return nil, err
	}

	// 2. Check if redirect URI matches
	if redirectURI != authorizationCode.RedirectURI.String {
		return nil, response.ErrInvalidRedirectURI
	}

	// 3. Check if the authorization code has expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, response.ErrAuthorizationCodeExpired
	}

	// 4. Log in the user
	accessToken, refreshToken, err := uc.Login(ctx, authorizationCode.Client, authorizationCode.User, authorizationCode.Scope)
	if err != nil {
		return nil, err
	}

	// 5. Delete the authorization code from the database
	err = uc.repository.DeleteAuthorizationCode(fmt.Sprint(authorizationCode.ID))
	if err != nil {
		return nil, err
	}

	// 6. Create the access token response
	accessTokenResponse, err := oauthDto.NewAccessTokenResponse(
		accessToken,
		refreshToken,
		config.BaseConfig.App.ConfigOauth.Oauth.AccessTokenLifetime,
		pkg.Bearer,
	)
	if err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}
