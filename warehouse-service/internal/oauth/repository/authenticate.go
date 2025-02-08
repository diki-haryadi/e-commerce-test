package oauthRepository

import (
	"context"
	"database/sql"
	"github.com/diki-haryadi/go-micro-template/config"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"time"
)

// Authenticate checks the access token is valid
func (rp *repository) Authenticate(token string) (*oauthDomain.AccessToken, error) {
	// 1. Fetch the access token from the database using a SELECT query
	sqlQuery := "SELECT id, token, client_id, user_id, expires_at FROM access_tokens WHERE token = $1"
	row := rp.postgres.SqlxDB.QueryRow(sqlQuery, token)

	// 2. Scan the results into an AccessToken object
	accessToken := new(oauthDomain.AccessToken)
	err := row.Scan(&accessToken.ID, &accessToken.Token, &accessToken.ClientID, &accessToken.UserID, &accessToken.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrAccessTokenNotFound
		}
		return nil, err
	}

	// 3. Check if the access token has expired
	if time.Now().UTC().After(accessToken.ExpiresAt) {
		return nil, response.ErrAccessTokenExpired
	}

	// 4. Extend the refresh token expiration
	// Use an UPDATE query to extend the expiration time for the refresh token
	refreshTokenQuery := `
        UPDATE refresh_tokens
        SET expires_at = $1
        WHERE client_id = $2 AND (user_id = $3 OR user_id IS NULL)
    `

	increasedExpiresAt := time.Now().UTC().Add(time.Duration(config.BaseConfig.App.ConfigOauth.Oauth.RefreshTokenLifetime) * time.Second)
	var userID interface{}
	if accessToken.UserID.Valid {
		userID = accessToken.UserID.String
	} else {
		userID = nil
	}

	// Execute the query to update the refresh token expiration
	_, err = rp.postgres.SqlxDB.Exec(refreshTokenQuery, increasedExpiresAt, accessToken.ClientID.String, userID)
	if err != nil {
		return nil, err
	}

	// Return the fetched access token
	return accessToken, nil
}

func (rp *repository) ClearUserTokens(userSession *oauthDomain.UserSession) {
	// 1. Check if the refresh token exists in the database
	tx, _ := rp.postgres.SqlxDB.BeginTx(context.Background(), nil)
	var refreshToken oauthDomain.RefreshToken
	sqlQuery := "SELECT * FROM refresh_tokens WHERE token = $1"
	row := tx.QueryRow(sqlQuery, userSession.RefreshToken)

	// 2. If refresh token is found, delete associated records with client_id and user_id
	err := row.Scan(&refreshToken.ID, &refreshToken.Token, &refreshToken.ClientID, &refreshToken.UserID)
	if err == nil { // Token found
		// Perform delete operation for refresh tokens
		deleteQuery := "DELETE FROM refresh_tokens WHERE client_id = $1 AND user_id = $2"
		_, err = tx.Exec(deleteQuery, refreshToken.ClientID.String, refreshToken.UserID.String)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	// 3. Check if the access token exists in the database
	var accessToken oauthDomain.AccessToken
	sqlQuery = "SELECT * FROM access_tokens WHERE token = $1"
	row = tx.QueryRow(sqlQuery, userSession.AccessToken)

	// 4. If access token is found, delete associated records with client_id and user_id
	err = row.Scan(&accessToken.ID, &accessToken.Token, &accessToken.ClientID, &accessToken.UserID)
	if err == nil { // Token found
		// Perform delete operation for access tokens
		deleteQuery := "DELETE FROM access_tokens WHERE client_id = $1 AND user_id = $2"
		_, err = tx.Exec(deleteQuery, accessToken.ClientID.String, accessToken.UserID.String)
		if err != nil {
			tx.Rollback()
			return
		}
	}
}
