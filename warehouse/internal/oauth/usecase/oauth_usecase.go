package oauthUseCase

import (
	sampleExtServiceDomain "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/domain"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
)

type useCase struct {
	repository              oauthDomain.Repository
	kafkaProducer           oauthDomain.KafkaProducer
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase
	allowedRoles            []string
}

func NewUseCase(
	repository oauthDomain.Repository,
	sampleExtServiceUseCase sampleExtServiceDomain.SampleExtServiceUseCase,
	kafkaProducer oauthDomain.KafkaProducer,
) oauthDomain.UseCase {
	return &useCase{
		repository:              repository,
		kafkaProducer:           kafkaProducer,
		sampleExtServiceUseCase: sampleExtServiceUseCase,
		allowedRoles:            []string{Superuser, User},
	}
}
