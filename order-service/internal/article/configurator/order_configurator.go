package articleConfigurator

import (
	"context"
	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	orderHttpController "github.com/diki-haryadi/go-micro-template/internal/article/delivery/http"
	orderKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/article/delivery/kafka/producer"
	orderDomain "github.com/diki-haryadi/go-micro-template/internal/article/domain"
	orderRepository "github.com/diki-haryadi/go-micro-template/internal/article/repository"
	orderUseCase "github.com/diki-haryadi/go-micro-template/internal/article/usecase"
	externalBridge "github.com/diki-haryadi/ztools/external_bridge"
	infraContainer "github.com/diki-haryadi/ztools/infra_container"
)

type configurator struct {
	ic        *infraContainer.IContainer
	extBridge *externalBridge.ExternalBridge
}

func NewConfigurator(ic *infraContainer.IContainer, extBridge *externalBridge.ExternalBridge) orderDomain.Configurator {
	return &configurator{ic: ic, extBridge: extBridge}
}

func (c *configurator) Configure(ctx context.Context) error {
	seServiceUseCase := sampleExtServiceUseCase.NewSampleExtServiceUseCase(c.extBridge.SampleExtGrpcService)
	kafkaProducer := orderKafkaProducer.NewProducer(c.ic.KafkaWriter)
	repository := orderRepository.NewRepository(c.ic.Postgres)
	useCase := orderUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// grpc
	//grpcController := orderGrpcController.NewController(useCase)
	//orderV1.RegisterArticleServiceServer(c.ic.GrpcServer.GetCurrentGrpcServer(), grpcController)

	// http
	httpRouterGp := c.ic.EchoHttpServer.GetEchoInstance().Group(c.ic.EchoHttpServer.GetBasePath())
	httpController := orderHttpController.NewController(useCase)
	orderHttpController.NewRouter(httpController).Register(httpRouterGp)

	// consumers
	//orderKafkaConsumer.NewConsumer(c.ic.KafkaReader).RunConsumers(ctx)

	// jobs
	//articleJob.NewJob(c.ic.Logger).StartJobs(ctx)

	return nil
}
