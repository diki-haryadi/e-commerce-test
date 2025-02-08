package oauthDomain

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
	Timestamp
}
