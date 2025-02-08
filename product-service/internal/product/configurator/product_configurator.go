package productConfigurator

import (
	"context"
	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	productHttpController "github.com/diki-haryadi/go-micro-template/internal/product/delivery/http"
	productKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/product/delivery/kafka/producer"
	productDomain "github.com/diki-haryadi/go-micro-template/internal/product/domain"
	productRepository "github.com/diki-haryadi/go-micro-template/internal/product/repository"
	productUseCase "github.com/diki-haryadi/go-micro-template/internal/product/usecase"
	externalBridge "github.com/diki-haryadi/ztools/external_bridge"
	infraContainer "github.com/diki-haryadi/ztools/infra_container"
)

type configurator struct {
	ic        *infraContainer.IContainer
	extBridge *externalBridge.ExternalBridge
}

func NewConfigurator(ic *infraContainer.IContainer, extBridge *externalBridge.ExternalBridge) productDomain.Configurator {
	return &configurator{ic: ic, extBridge: extBridge}
}

func (c *configurator) Configure(ctx context.Context) error {
	seServiceUseCase := sampleExtServiceUseCase.NewSampleExtServiceUseCase(c.extBridge.SampleExtGrpcService)
	kafkaProducer := productKafkaProducer.NewProducer(c.ic.KafkaWriter)
	repository := productRepository.NewRepository(c.ic.Postgres)
	useCase := productUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// grpc
	//grpcController := productGrpcController.NewController(useCase)
	//productV1.RegisterArticleServiceServer(c.ic.GrpcServer.GetCurrentGrpcServer(), grpcController)

	// http
	httpRouterGp := c.ic.EchoHttpServer.GetEchoInstance().Group(c.ic.EchoHttpServer.GetBasePath())
	httpController := productHttpController.NewController(useCase)
	productHttpController.NewRouter(httpController).Register(httpRouterGp)

	// consumers
	//productKafkaConsumer.NewConsumer(c.ic.KafkaReader).RunConsumers(ctx)

	// jobs
	//productJob.NewJob(c.ic.Logger).StartJobs(ctx)

	return nil
}
