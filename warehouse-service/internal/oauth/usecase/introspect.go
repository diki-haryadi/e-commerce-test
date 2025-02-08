package oauthUseCase

import (
	"context"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/constant"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
)

func (uc *useCase) IntrospectToken(ctx context.Context, token, tokenTypeHint string, client *oauthDomain.Client) (*oauthDto.IntrospectResponse, error) {
	// Default to access token hint
	if tokenTypeHint == "" {
		tokenTypeHint = constant.AccessTokenHint
	}

	switch tokenTypeHint {
	case constant.AccessTokenHint:
		accessToken, err := uc.repository.Authenticate(token)
		if err != nil {
			return nil, err
		}
		return uc.NewIntrospectResponseFromAccessToken(accessToken)
	case constant.RefreshTokenHint:
		refreshToken, err := uc.repository.GetValidRefreshToken(token, client)
		if err != nil {
			return nil, err
		}
		return uc.NewIntrospectResponseFromRefreshToken(refreshToken)
	default:
		return nil, response.ErrTokenHintInvalid
	}
}

// NewIntrospectResponseFromAccessToken creates an introspect.md response from an access token
func (uc *useCase) NewIntrospectResponseFromAccessToken(accessToken *oauthDomain.AccessToken) (*oauthDto.IntrospectResponse, error) {
	var introspectResponse = &oauthDto.IntrospectResponse{
		Active:    true,
		Scope:     accessToken.Scope,
		TokenType: pkg.Bearer,
		ExpiresAt: int(accessToken.ExpiresAt.Unix()),
	}
	// Fetch the client using the FetchClientByClientID method
	if accessToken.ClientID.Valid {
		client, err := uc.repository.FetchClientByClientID(accessToken.ClientID.String)
		if err != nil {
			return nil, err
		}
		introspectResponse.ClientID = client.Key
	}

	// Fetch the user using the FetchUserByUserID method
	if accessToken.UserID.Valid {
		user, err := uc.repository.FetchUserByUserID(accessToken.UserID.String)
		if err != nil {
			return nil, err
		}
		introspectResponse.Username = user.Username
	}

	return introspectResponse, nil
}

// NewIntrospectResponseFromRefreshToken creates an introspect.md response from a refresh token
func (uc *useCase) NewIntrospectResponseFromRefreshToken(refreshToken *oauthDomain.RefreshToken) (*oauthDto.IntrospectResponse, error) {
	var introspectResponse = &oauthDto.IntrospectResponse{
		Active:    true,
		Scope:     refreshToken.Scope,
		TokenType: pkg.Bearer,
		ExpiresAt: int(refreshToken.ExpiresAt.Unix()),
	}

	// Fetch the client using the FetchClientByClientID method
	if refreshToken.ClientID.Valid {
		client, err := uc.repository.FetchClientByClientID(refreshToken.ClientID.String)
		if err != nil {
			return nil, err
		}
		introspectResponse.ClientID = client.Key
	}

	// Fetch the user using the FetchUserByUserID method
	if refreshToken.UserID.Valid {
		user, err := uc.repository.FetchUserByUserID(refreshToken.UserID.String)
		if err != nil {
			return nil, err
		}
		introspectResponse.Username = user.Username
	}

	return introspectResponse, nil
}
