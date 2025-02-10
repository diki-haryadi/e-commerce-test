package articleRepository

import (
	"context"
	"fmt"

	authDomain "github.com/diki-haryadi/go-micro-template/internal/auth/domain"
	authDto "github.com/diki-haryadi/go-micro-template/internal/auth/dto"
	"github.com/diki-haryadi/ztools/postgres"
)

type repository struct {
	postgres *postgres.Postgres
}

func NewRepository(conn *postgres.Postgres) authDomain.Repository {
	return &repository{postgres: conn}
}

func (rp *repository) SignUp(ctx context.Context, entity *authDto.SignUpRequestDto) (*authDto.CreateSignUpResponseDto, error) {
	query := `INSERT INTO public.users (username, password) VALUES ($1, $2) RETURNING id, username`

	result, err := rp.postgres.SqlxDB.QueryContext(ctx, query, entity.Username, entity.Password)
	if err != nil {
		return nil, fmt.Errorf("error inserting auth record")
	}

	user := new(authDto.CreateSignUpResponseDto)
	for result.Next() {
		err = result.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (rp *repository) GetUserByUsername(ctx context.Context, username string) (*authDto.CreateSignInResponseDto, error) {
	query := `SELECT id, username, password FROM public.users where username = $1`

	user := new(authDto.CreateSignInResponseDto)
	err := rp.postgres.SqlxDB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (rp *repository) GetUserById(ctx context.Context, userId string) (*authDto.ProfileResponse, error) {
	query := `
        SELECT id, username
        FROM users 
        WHERE id = $1
    `

	profile := new(authDto.ProfileResponse)
	err := rp.postgres.SqlxDB.QueryRowContext(ctx, query, userId).Scan(
		&profile.ID,
		&profile.Username,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile: %w", err)
	}

	return profile, nil
}
