package articleDomain

import (
	"context"
	"github.com/google/uuid"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"

	authDto "github.com/diki-haryadi/go-micro-template/internal/auth/dto"
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
	SignUp(ctx context.Context, article *authDto.SignUpRequestDto) (*authDto.CreateSignUpResponseDto, error)
	SignIn(ctx context.Context, username, password string) (*authDto.SignInResponse, error)
	GetProfile(ctx context.Context, userId string) (*authDto.ProfileResponse, error)
}

type Repository interface {
	SignUp(ctx context.Context, article *authDto.SignUpRequestDto) (*authDto.CreateSignUpResponseDto, error)
	GetUserByUsername(ctx context.Context, username string) (*authDto.CreateSignInResponseDto, error)
	GetUserById(ctx context.Context, userId string) (*authDto.ProfileResponse, error)
}

type GrpcController interface {
	//CreateArticle(ctx context.Context, req *articleV1.CreateArticleRequest) (*articleV1.CreateArticleResponse, error)
	//GetArticleById(ctx context.Context, req *articleV1.GetArticleByIdRequest) (*articleV1.GetArticleByIdResponse, error)
}

type HttpController interface {
	SignUp(c echo.Context) error
	SignIn(ctx echo.Context) error
	RefreshToken(ctx echo.Context) error
	GetProfile(ctx echo.Context) error
}

type Job interface {
	//StartJobs(ctx context.Context)
}

type KafkaProducer interface {
	PublishCreateEvent(ctx context.Context, messages ...kafka.Message) error
}

type KafkaConsumer interface {
	RunConsumers(ctx context.Context)
}
