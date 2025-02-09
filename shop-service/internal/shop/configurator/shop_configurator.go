package shopConfigurator

import (
	"context"
	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	shopHttpController "github.com/diki-haryadi/go-micro-template/internal/shop/delivery/http"
	shopKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/shop/delivery/kafka/producer"
	shopDomain "github.com/diki-haryadi/go-micro-template/internal/shop/domain"
	shopRepository "github.com/diki-haryadi/go-micro-template/internal/shop/repository"
	shopUseCase "github.com/diki-haryadi/go-micro-template/internal/shop/usecase"
	externalBridge "github.com/diki-haryadi/ztools/external_bridge"
	infraContainer "github.com/diki-haryadi/ztools/infra_container"
)

type configurator struct {
	ic        *infraContainer.IContainer
	extBridge *externalBridge.ExternalBridge
}

func NewConfigurator(ic *infraContainer.IContainer, extBridge *externalBridge.ExternalBridge) shopDomain.Configurator {
	return &configurator{ic: ic, extBridge: extBridge}
}

func (c *configurator) Configure(ctx context.Context) error {
	seServiceUseCase := sampleExtServiceUseCase.NewSampleExtServiceUseCase(c.extBridge.SampleExtGrpcService)
	kafkaProducer := shopKafkaProducer.NewProducer(c.ic.KafkaWriter)
	repository := shopRepository.NewRepository(c.ic.Postgres)
	useCase := shopUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// grpc
	//grpcController := shopGrpcController.NewController(useCase)
	//shopV1.RegisterArticleServiceServer(c.ic.GrpcServer.GetCurrentGrpcServer(), grpcController)

	// http
	httpRouterGp := c.ic.EchoHttpServer.GetEchoInstance().Group(c.ic.EchoHttpServer.GetBasePath())
	httpController := shopHttpController.NewController(useCase)
	shopHttpController.NewRouter(httpController).Register(httpRouterGp)

	// consumers
	//shopKafkaConsumer.NewConsumer(c.ic.KafkaReader).RunConsumers(ctx)

	// jobs
	//shopJob.NewJob(c.ic.Logger).StartJobs(ctx)

	return nil
}
