package oauthDomain

import (
	"database/sql"
	"time"
)

type AuthorizationCode struct {
	Common
	ClientID    sql.NullString `db:"client_id"`
	UserID      sql.NullString `db:"user_id"`
	Client      *Client
	User        *Users
	Code        string         `sql:"code"`
	RedirectURI sql.NullString `db:"redirect_uri"`
	ExpiresAt   time.Time      `sql:"expires_at"`
	Scope       string         `sql:"scope"`
}
