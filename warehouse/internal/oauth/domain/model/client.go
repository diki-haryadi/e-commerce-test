package oauthDomain

import "database/sql"

type Client struct {
	Common
	Key         string         `db:"key" json:"key"`
	Secret      string         `db:"secret" json:"secret"`
	RedirectURI sql.NullString `db:"redirect_uri" json:"redirect_uri"`
}
