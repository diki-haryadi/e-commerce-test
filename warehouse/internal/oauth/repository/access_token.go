package oauthRepository

import (
	"context"
	"fmt"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	"time"
)

func (rp *repository) GrantAccessToken(ctx context.Context, client *oauthDomain.Client, user *oauthDomain.Users, expiresIn int, scope string) (*oauthDomain.AccessToken, error) {
	// Begin a transaction
	tx, err := rp.postgres.SqlxDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	rawQuery := " WHERE client_id = $1"
	args := []interface{}{fmt.Sprint(client.ID)}

	// Add the user_id condition if necessary
	if user != nil && fmt.Sprint(user) != "" {
		rawQuery += " AND user_id = $2"
		args = append(args, user.ID) // Add user.ID to the arguments list
	} else {
		rawQuery += " AND user_id IS NULL"
	}

	// Add the expiration condition
	rawQuery += " AND expires_at <= NOW()"

	// Complete the DELETE statement
	query := "DELETE FROM access_tokens" + rawQuery

	// Execute the query using parameterized arguments
	_, err = tx.Exec(query, args...)
	if err != nil {
		// If an error occurs, rollback the transaction
		_ = tx.Rollback()
		return nil, err
	}

	// Create a new access token
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)
	accessToken := oauthDto.NewOauthAccessToken(client, user, expiresIn, scope)

	var sqlQueryAT string
	var insertArgs []interface{}

	if user != nil {
		sqlQueryAT = `
        INSERT INTO access_tokens (client_id, user_id, token, expires_at, scope)
        VALUES ($1, $2, $3, $4, $5)`
		insertArgs = append(insertArgs, fmt.Sprint(client.ID), user.ID)
	} else {
		sqlQueryAT = `
        INSERT INTO access_tokens (client_id, token, expires_at, scope)
        VALUES ($1, $2, $3, $4)`
		insertArgs = append(insertArgs, fmt.Sprint(client.ID))
	}

	insertArgs = append(insertArgs, accessToken.Token, expiresAt, scope)
	_, err = tx.Exec(sqlQueryAT, insertArgs...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	accessToken.Client = client
	accessToken.User = user

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return accessToken, nil
}
