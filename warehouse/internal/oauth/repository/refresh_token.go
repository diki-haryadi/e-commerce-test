package oauthRepository

import (
	"database/sql"
	"fmt"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"time"
)

// GetOrCreateRefreshToken retrieves an existing refresh token, if expired,
// the token gets deleted and a new refresh token is created using raw SQL
func (rp *repository) GetOrCreateRefreshToken(client *oauthDomain.Client, user *oauthDomain.Users, expiresIn int, scope string) (*oauthDomain.RefreshToken, error) {
	// Try to fetch an existing refresh token first using raw SQL
	var refreshToken oauthDomain.RefreshToken
	sqlQuery := "SELECT id, client_id, user_id, token, expires_at, scope FROM refresh_tokens WHERE client_id = $1"

	if user != nil && fmt.Sprint(user.ID) != "" {
		sqlQuery += " AND user_id = $2"
	} else {
		sqlQuery += " AND user_id IS NULL"
	}

	err := rp.postgres.SqlxDB.QueryRow(sqlQuery, client.ID, user.ID).Scan(
		&refreshToken.ID,
		&refreshToken.ClientID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.Scope,
	)

	var expired bool
	if err == nil {
		// Check if the token is expired
		expired = time.Now().UTC().After(refreshToken.ExpiresAt)
	}

	// If the refresh token has expired or does not exist, delete it
	if expired || err == sql.ErrNoRows {
		if err == nil { // If token exists, delete it
			sqlDelete := "DELETE FROM refresh_tokens WHERE id = $1"
			_, err = rp.postgres.SqlxDB.Exec(sqlDelete, refreshToken.ID)
			if err != nil {
				return nil, err
			}
		}

		// Create a new refresh token if it expired or was not found
		refreshToken = *oauthDto.NewOauthRefreshToken(client, user, expiresIn, scope)

		sqlInsert := `
            INSERT INTO refresh_tokens (client_id, user_id, token, expires_at, scope)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, client_id, user_id, token, expires_at, scope`
		fmt.Println(refreshToken)
		err = rp.postgres.SqlxDB.QueryRow(sqlInsert, refreshToken.ClientID, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, refreshToken.Scope).
			Scan(&refreshToken.ID, &refreshToken.ClientID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.Scope)

		if err != nil {
			return nil, err
		}
	}

	return &refreshToken, nil
}

// GetValidRefreshToken returns a valid non expired refresh token using raw SQL
func (rp *repository) GetValidRefreshToken(token string, client *oauthDomain.Client) (*oauthDomain.RefreshToken, error) {
	// Fetch the refresh token from the database using raw SQL
	var refreshToken oauthDomain.RefreshToken
	sqlQuery := "SELECT id, client_id, user_id, token, expires_at, scope FROM refresh_tokens WHERE client_id = $1 AND token = $2"

	err := rp.postgres.SqlxDB.QueryRow(sqlQuery, client.ID, token).Scan(
		&refreshToken.ID,
		&refreshToken.ClientID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.Scope,
	)

	// Not found
	if err == sql.ErrNoRows {
		return nil, response.ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	// Check if the refresh token hasn't expired
	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		return nil, response.ErrRefreshTokenExpired
	}

	return &refreshToken, nil
}
