package oauthDomain

import "database/sql"

type Users struct {
	Common
	RoleID   sql.NullString `db:"role_id" json:"role_id"`
	Role     *Role
	Username string         `db:"username" json:"username"`
	Password sql.NullString `db:"password" json:"password"`
}
