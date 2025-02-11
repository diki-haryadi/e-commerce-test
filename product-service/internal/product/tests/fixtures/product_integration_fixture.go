package productFixture

import (
	"context"
	"math"
	"time"

	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	productHttp "github.com/diki-haryadi/go-micro-template/internal/product/delivery/http"
	productKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/product/delivery/kafka/producer"
	productRepo "github.com/diki-haryadi/go-micro-template/internal/product/repository"
	productUseCase "github.com/diki-haryadi/go-micro-template/internal/product/usecase"
	externalBridge "github.com/diki-haryadi/ztools/external_bridge"
	iContainer "github.com/diki-haryadi/ztools/infra_container"
)

type IntegrationTestFixture struct {
	TearDown       func()
	Ctx            context.Context
	Cancel         context.CancelFunc
	InfraContainer *iContainer.IContainer
}

func NewIntegrationTestFixture() (*IntegrationTestFixture, error) {
	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	container := iContainer.IContainer{}
	ic, infraDown, err := container.IContext(ctx).
		ICDown().ICPostgres().ICEcho().NewIC()
	if err != nil {
		cancel()
		return nil, err
	}

	extBridge, extBridgeDown, err := externalBridge.NewExternalBridge(ctx)
	if err != nil {
		cancel()
		return nil, err
	}

	seServiceUseCase := sampleExtServiceUseCase.NewSampleExtServiceUseCase(extBridge.SampleExtGrpcService)
	kafkaProducer := productKafkaProducer.NewProducer(ic.KafkaWriter)
	repository := productRepo.NewRepository(ic.Postgres)
	useCase := productUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// http
	ic.EchoHttpServer.SetupDefaultMiddlewares()
	httpRouterGp := ic.EchoHttpServer.GetEchoInstance().Group(ic.EchoHttpServer.GetBasePath())
	httpController := productHttp.NewController(useCase)
	productHttp.NewRouter(httpController).Register(httpRouterGp)

	return &IntegrationTestFixture{
		TearDown: func() {
			cancel()
			infraDown()
			extBridgeDown()
		},
		InfraContainer: ic,
		Ctx:            ctx,
		Cancel:         cancel,
	}, nil
}
