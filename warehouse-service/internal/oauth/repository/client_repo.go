package oauthRepository

import (
	"context"
	"database/sql"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/google/uuid"
	"strings"
	"time"
)

func (rp *repository) CreateClientCommon(ctx context.Context, clientID, secret, redirectURI string) (*oauthDomain.Client, error) {
	// 1. Check if client already exists
	var existingClient oauthDomain.Client
	sqlCheck := `SELECT id FROM clients WHERE client_id = $1`
	err := rp.postgres.SqlxDB.Get(&existingClient, sqlCheck, clientID)
	if err == nil {
		return nil, response.ErrClientIDTaken // Client ID is already taken
	}
	if err != sql.ErrNoRows {
		return nil, err // Other errors
	}

	// 2. Hash the secret (password)
	secretHash, err := pkg.HashPassword(secret)
	if err != nil {
		return nil, err
	}

	// 3. Insert the new client into the database
	sqlInsert := `
        INSERT INTO clients (client_id, secret, redirect_uri, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, client_id, secret, redirect_uri, created_at
    `

	client := &oauthDomain.Client{
		Common: oauthDomain.Common{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Key:         strings.ToLower(clientID),
		Secret:      string(secretHash),
		RedirectURI: pkg.StringOrNull(redirectURI),
	}

	// Execute the insert query and scan the results into the client struct
	err = rp.postgres.SqlxDB.QueryRow(sqlInsert, client.Key, client.Secret, client.RedirectURI, client.CreatedAt).Scan(
		&client.ID, &client.Key, &client.Secret, &client.RedirectURI, &client.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (rp *repository) FindClientByClientID(ctx context.Context, clientID string) (*oauthDomain.Client, error) {
	client := oauthDomain.Client{}
	query := "SELECT * FROM clients WHERE key = $1"
	err := rp.postgres.SqlxDB.Get(&client, query, strings.ToLower(clientID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrClientNotFound
		}
		return nil, err
	}

	return &client, err
}
