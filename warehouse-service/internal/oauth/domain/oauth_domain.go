package oauthDomain

import (
	"context"
	model "github.com/diki-haryadi/go-micro-template/internal/oauth/domain/model"
	oauthDto "github.com/diki-haryadi/go-micro-template/internal/oauth/dto"
	articleV1 "github.com/diki-haryadi/protobuf-template/go-micro-template/article/v1"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

type Article struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"desc"`
}

type Configurator interface {
	Configure(ctx context.Context) error
}

type UseCase interface {
	AuthClient(ctx context.Context, clientID, secret string) (*model.Client, error)
	CreateClient(ctx context.Context, clientID, secret, redirectURI string) (*model.Client, error)
	ClientExists(ctx context.Context, clientID string) bool
	ClientCredentialsGrant(ctx context.Context, scope string, refreshToken string, client *model.Client) (*oauthDto.AccessTokenResponse, error)
	AuthorizationCodeGrant(ctx context.Context, code, redirectURI string, client *model.Client) (*oauthDto.AccessTokenResponse, error)
	PasswordGrant(ctx context.Context, username, password string, scope string, client *model.Client) (*oauthDto.AccessTokenResponse, error)
	RefreshTokenGrant(ctx context.Context, token, scope string, client *model.Client) (*oauthDto.AccessTokenResponse, error)
	IntrospectToken(ctx context.Context, token, tokenTypeHint string, client *model.Client) (*oauthDto.IntrospectResponse, error)
	Login(ctx context.Context, client *model.Client, user *model.Users, scope string) (*model.AccessToken, *model.RefreshToken, error)
	GetRefreshTokenScope(ctx context.Context, refreshToken *model.RefreshToken, requestedScope string) (string, error)
	GrantAccessToken(ctx context.Context, client *model.Client, user *model.Users, expiresIn int, scope string) (*oauthDto.AccessTokenResponse, error)
}

type Repository interface {
	GrantAccessToken(ctx context.Context, client *model.Client, user *model.Users, expiresIn int, scope string) (*model.AccessToken, error)
	Authenticate(token string) (*model.AccessToken, error)
	GrantAuthorizationCode(client *model.Client, user *model.Users, expiresIn int, redirectURI, scope string) (*model.AuthorizationCode, error)
	GetValidAuthorizationCode(code, redirectURI string, client *model.Client) (*model.AuthorizationCode, error)
	CreateClientCommon(ctx context.Context, clientID, secret, redirectURI string) (*model.Client, error)
	FindClientByClientID(ctx context.Context, clientID string) (*model.Client, error)
	FetchAuthorizationCodeByCode(ctx context.Context, client *model.Client, code string) (*model.AuthorizationCode, error)
	DeleteAuthorizationCode(authorizationCodeID string) error
	FetchClientByClientID(clientID string) (*model.Client, error)
	FetchUserByUserID(userID string) (*model.Users, error)
	GetOrCreateRefreshToken(client *model.Client, user *model.Users, expiresIn int, scope string) (*model.RefreshToken, error)
	GetValidRefreshToken(token string, client *model.Client) (*model.RefreshToken, error)
	FindRoleByID(id string) (*model.Role, error)
	GetScope(ctx context.Context, requestedScope string) (string, error)
	GetDefaultScope(ctx context.Context) string
	ScopeExists(ctx context.Context, requestedScope string) bool
	FindUserByUsername(username string) (*model.Users, error)
	CreateUserCommon(roleID, username, password string) (*model.Users, error)
	SetPasswordCommon(user *model.Users, password string) error
	UpdateUsernameCommon(user *model.Users, username string) error
}

type GrpcController interface {
	CreateArticle(ctx context.Context, req *articleV1.CreateArticleRequest) (*articleV1.CreateArticleResponse, error)
	GetArticleById(ctx context.Context, req *articleV1.GetArticleByIdRequest) (*articleV1.GetArticleByIdResponse, error)
}

type HttpController interface {
	Tokens(c echo.Context) error
	Introspect(c echo.Context) error
}

type Job interface {
	StartJobs(ctx context.Context)
}

type KafkaProducer interface {
	PublishCreateEvent(ctx context.Context, messages ...kafka.Message) error
}

type KafkaConsumer interface {
	RunConsumers(ctx context.Context)
}
