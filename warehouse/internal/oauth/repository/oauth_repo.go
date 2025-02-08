package oauthRepository

import (
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
	"github.com/diki-haryadi/ztools/postgres"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepository(conn *postgres.Postgres) oauthDomain.Repository {
	return &repository{postgres: conn}
}
