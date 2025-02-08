package oauthDomain

import (
	"github.com/google/uuid"
	"time"
)

type Common struct {
	ID        uuid.UUID  `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Timestamp struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type EmailTokenModel struct {
	Common
	Reference   string     `db:"reference"`
	EmailSent   bool       `db:"email_sent"`
	EmailSentAt *time.Time `db:"email_sent_at"`
	ExpiresAt   time.Time  `db:"expires_at"`
}
