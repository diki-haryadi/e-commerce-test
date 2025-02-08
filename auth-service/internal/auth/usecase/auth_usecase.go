package articleUseCase

import (
	"context"
	"encoding/json"
	"github.com/diki-haryadi/go-micro-template/config"
	"github.com/diki-haryadi/go-micro-template/pkg"

	"log"

	"github.com/segmentio/kafka-go"

	sampleExtServiceDomain "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/domain"
	authDomain "github.com/diki-haryadi/go-micro-template/internal/auth/domain"
	authDto "github.com/diki-haryadi/go-micro-template/internal/auth/dto"
)

type useCase struct {
	repository              authDomain.Repository
	kafkaProducer           authDomain.KafkaProducer
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase
}

func NewUseCase(
	repository authDomain.Repository,
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase,
	kafkaProducer authDomain.KafkaProducer,
) authDomain.UseCase {
	return &useCase{
		repository:              repository,
		kafkaProducer:           kafkaProducer,
		sampleExtServiceUseCase: sampleExtServiceUseCase,
	}
}

func (uc *useCase) SignUp(ctx context.Context, req *authDto.SignUpRequestDto) (*authDto.CreateSignUpResponseDto, error) {
	user, err := uc.repository.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}

	// TODO : if err => return Marshal_Err_Exception
	jsonArticle, _ := json.Marshal(user)

	// if it has go keyword and if we pass the request context to it, it will terminate after request lifecycle.
	_ = uc.kafkaProducer.PublishCreateEvent(context.Background(), kafka.Message{
		Key:   []byte("Article"),
		Value: jsonArticle,
	})

	return user, err
}

func (uc *useCase) SignIn(ctx context.Context, username, password string) (*authDto.SignInResponse, error) {
	resp := &authDto.SignInResponse{}
	user, err := uc.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if !user.Password.Valid {
		return nil, err
	}

	if pkg.VerifyPassword(user.Password.String, password) != nil {
		return nil, err
	}

	tokenPair, err := pkg.GenerateTokenPair(
		user.ID.String(),
		"auth-service", // issuer
		[]string{"auth-service", "product-service", "order-service", "warehouse-service"}, // audience
		"user:read user:write", // scope
		config.BaseConfig.App.JWTSecret,
	)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
		return nil, err
	}

	resp.User = *user
	resp.AccessToken = tokenPair.AccessToken
	resp.RefreshToken = tokenPair.RefreshToken
	resp.ExpiresIn = tokenPair.ExpiresIn
	return resp, nil
}

func (uc *useCase) GetProfile(ctx context.Context, userId string) (*authDto.ProfileResponse, error) {
	profile, err := uc.repository.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
