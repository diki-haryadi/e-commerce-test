package oauthDomain

import (
	"database/sql"
	"time"
)

type Token struct {
	Common
	ClientID  sql.NullString `db:"client_id"`
	UserID    sql.NullString `db:"user_id"`
	Client    *Client
	User      *Users
	Token     string    `sql:"token"`
	ExpiresAt time.Time `sql:"expires_at"`
	Scope     string    `sql:"scope"`
}
