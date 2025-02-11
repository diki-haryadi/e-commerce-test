package productFixture

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/app"
	"math"
	"time"

	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	whHttp "github.com/diki-haryadi/go-micro-template/internal/warehouse/delivery/http"
	whKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/warehouse/delivery/kafka/producer"
	whRepo "github.com/diki-haryadi/go-micro-template/internal/warehouse/repository"
	whUseCase "github.com/diki-haryadi/go-micro-template/internal/warehouse/usecase"
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
	_ = app.New().Init()
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
	kafkaProducer := whKafkaProducer.NewProducer(ic.KafkaWriter)
	repository := whRepo.NewRepository(ic.Postgres)
	useCase := whUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// http
	ic.EchoHttpServer.SetupDefaultMiddlewares()
	httpRouterGp := ic.EchoHttpServer.GetEchoInstance().Group(ic.EchoHttpServer.GetBasePath())
	httpController := whHttp.NewController(useCase)
	whHttp.NewRouter(httpController).Register(httpRouterGp)

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
