package oauthRepository

import (
	"database/sql"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"time"
)

// GrantAuthorizationCode grants a new authorization code using raw SQL
func (rp *repository) GrantAuthorizationCode(client *oauthDomain.Client, user *oauthDomain.Users, expiresIn int, redirectURI, scope string) (*oauthDomain.AuthorizationCode, error) {
	// Generate a new authorization code
	authorizationCode := oauthDto.NewOauthAuthorizationCode(client, user, expiresIn, redirectURI, scope)

	// Prepare the SQL INSERT query to insert the new authorization code into the database
	sqlQuery := `
        INSERT INTO authorization_codes (client_id, user_id, code, redirect_uri, expires_at, scope)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, client_id, user_id, code, redirect_uri, expires_at, scope
    `

	// Execute the query and retrieve the generated ID and other fields
	row := rp.postgres.SqlxDB.QueryRow(sqlQuery, client.ID, user.ID, authorizationCode.Code, authorizationCode.RedirectURI.String, authorizationCode.ExpiresAt, authorizationCode.Scope)

	// Map the result into the authorizationCode object
	err := row.Scan(&authorizationCode.ID, &authorizationCode.ClientID, &authorizationCode.UserID, &authorizationCode.Code, &authorizationCode.RedirectURI, &authorizationCode.ExpiresAt, &authorizationCode.Scope)
	if err != nil {
		return nil, err
	}

	// Set the associated client and user (these are already set from the input)
	authorizationCode.Client = client
	authorizationCode.User = user

	return authorizationCode, nil
}

// getValidAuthorizationCode returns a valid non-expired authorization code using raw SQL
func (rp *repository) GetValidAuthorizationCode(code, redirectURI string, client *oauthDomain.Client) (*oauthDomain.AuthorizationCode, error) {
	// Fetch the authorization code from the database using raw SQL query
	sqlQuery := `
        SELECT id, client_id, user_id, code, redirect_uri, expires_at, scope
        FROM authorization_codes
        WHERE client_id = $1 AND code = $2
    `

	row := rp.postgres.SqlxDB.QueryRow(sqlQuery, client.ID, code)

	// Scan the result into an authorizationCode object
	authorizationCode := new(oauthDomain.AuthorizationCode)
	err := row.Scan(&authorizationCode.ID, &authorizationCode.ClientID, &authorizationCode.UserID, &authorizationCode.Code, &authorizationCode.RedirectURI, &authorizationCode.ExpiresAt, &authorizationCode.Scope)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrAuthorizationCodeNotFound
		}
		return nil, err
	}

	// Check if the redirect URI matches
	if redirectURI != authorizationCode.RedirectURI.String {
		return nil, response.ErrInvalidRedirectURI
	}

	// Check if the authorization code has expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, response.ErrAuthorizationCodeExpired
	}

	return authorizationCode, nil
}
