package oauthUseCase

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
)

func (uc *useCase) AuthClient(ctx context.Context, clientID, secret string) (*oauthDomain.Client, error) {
	client, err := uc.repository.FindClientByClientID(ctx, clientID)
	if err != nil {
		return nil, response.ErrClientNotFound
	}

	if pkg.VerifyPassword(client.Secret, secret) != nil {
		return nil, response.ErrInvalidClientSecret
	}
	return client, nil
}

func (uc *useCase) CreateClient(ctx context.Context, clientID, secret, redirectURI string) (*oauthDomain.Client, error) {
	client, err := uc.repository.CreateClientCommon(ctx, clientID, secret, redirectURI)
	if err != nil {
		return nil, err
	}
	return client, err
}

func (uc *useCase) ClientExists(ctx context.Context, clientID string) bool {
	_, err := uc.repository.FindClientByClientID(ctx, clientID)
	return err == nil
}
