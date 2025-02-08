package oauthConfigurator

import (
	"context"
	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	oauthHttpController "github.com/diki-haryadi/go-micro-template/internal/oauth/delivery/http"
	oauthKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/oauth/delivery/kafka/producer"
	oauthDomain "github.com/diki-haryadi/go-micro-template/internal/oauth/domain"
	oauthRepository "github.com/diki-haryadi/go-micro-template/internal/oauth/repository"
	oauthUseCase "github.com/diki-haryadi/go-micro-template/internal/oauth/usecase"
	externalBridge "github.com/diki-haryadi/ztools/external_bridge"
	infraContainer "github.com/diki-haryadi/ztools/infra_container"
)

type configurator struct {
	ic        *infraContainer.IContainer
	extBridge *externalBridge.ExternalBridge
}

func NewConfigurator(ic *infraContainer.IContainer, extBridge *externalBridge.ExternalBridge) oauthDomain.Configurator {
	return &configurator{ic: ic, extBridge: extBridge}
}

func (c *configurator) Configure(ctx context.Context) error {
	seServiceUseCase := sampleExtServiceUseCase.NewSampleExtServiceUseCase(c.extBridge.SampleExtGrpcService)
	kafkaProducer := oauthKafkaProducer.NewProducer(c.ic.KafkaWriter)
	repository := oauthRepository.NewRepository(c.ic.Postgres)
	useCase := oauthUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// grpc
	//grpcController := oauthGrpcController.NewController(useCase)
	//articleV1.RegisterArticleServiceServer(c.ic.GrpcServer.GetCurrentGrpcServer(), grpcController)

	// http
	httpRouterGp := c.ic.EchoHttpServer.GetEchoInstance().Group(c.ic.EchoHttpServer.GetBasePath())
	httpController := oauthHttpController.NewController(useCase)
	oauthHttpController.NewRouter(httpController).Register(httpRouterGp)

	// consumers
	//oauthKafkaConsumer.NewConsumer(c.ic.KafkaReader).RunConsumers(ctx)

	// jobs
	//oauthJob.NewJob(c.ic.Logger).StartJobs(ctx)

	return nil
}
