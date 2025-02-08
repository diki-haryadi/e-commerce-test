package oauthRepository

import (
	"context"
	"database/sql"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	"github.com/diki-haryadi/go-micro-template/pkg/response"
	"github.com/google/uuid"
)

// FetchAuthorizationCodeByCode retrieves the authorization code from the database using raw SQL
func (rp *repository) FetchAuthorizationCodeByCode(ctx context.Context, client *oauthDomain.Client, code string) (*oauthDomain.AuthorizationCode, error) {
	sqlQuery := `
        SELECT ac.id, ac.client_id, ac.user_id, ac.code, ac.redirect_uri, ac.expires_at, ac.scope, r.name AS role_name
		FROM authorization_codes ac
		JOIN users u ON ac.user_id::UUID = u.id
		JOIN roles r ON u.role_id::UUID = r.id
        WHERE client_id = $1 AND code = $2`

	var authorizationCode oauthDomain.AuthorizationCode
	var user oauthDomain.Users
	var role oauthDomain.Role
	var cl oauthDomain.Client

	row := rp.postgres.SqlxDB.QueryRow(sqlQuery, client.ID, code)

	// Scan the result into the authorizationCode struct
	err := row.Scan(
		&authorizationCode.ID,
		&authorizationCode.ClientID,
		&authorizationCode.UserID,
		&authorizationCode.Code,
		&authorizationCode.RedirectURI,
		&authorizationCode.ExpiresAt,
		&authorizationCode.Scope,
		&role.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrAuthorizationCodeNotFound
		}
		return nil, err
	}

	cl.ID = client.ID
	u, _ := uuid.Parse(authorizationCode.UserID.String)
	user.ID = u
	user.Role = &role

	authorizationCode.User = &user
	authorizationCode.Client = &cl

	return &authorizationCode, nil
}

// DeleteAuthorizationCode deletes the authorization code from the database after use
func (rp *repository) DeleteAuthorizationCode(authorizationCodeID string) error {
	sqlDelete := `
        DELETE FROM authorization_codes WHERE id = $1
    `
	_, err := rp.postgres.SqlxDB.Exec(sqlDelete, authorizationCodeID)
	return err
}
