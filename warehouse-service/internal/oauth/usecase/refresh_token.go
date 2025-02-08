package oauthUseCase

import (
	"context"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
)

// GetRefreshTokenScope returns scope for a new refresh token
func (uc *useCase) GetRefreshTokenScope(ctx context.Context, refreshToken *oauthDomain.RefreshToken, requestedScope string) (string, error) {
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
