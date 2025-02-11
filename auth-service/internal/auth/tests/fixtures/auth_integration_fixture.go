package authFixture

import (
	"context"
	"github.com/diki-haryadi/go-micro-template/app"
	"math"
	"time"

	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	authHttp "github.com/diki-haryadi/go-micro-template/internal/auth/delivery/http"
	authKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/auth/delivery/kafka/producer"
	authRepo "github.com/diki-haryadi/go-micro-template/internal/auth/repository"
	authUseCase "github.com/diki-haryadi/go-micro-template/internal/auth/usecase"
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
	kafkaProducer := authKafkaProducer.NewProducer(ic.KafkaWriter)
	repository := authRepo.NewRepository(ic.Postgres)
	useCase := authUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// http
	ic.EchoHttpServer.SetupDefaultMiddlewares()
	httpRouterGp := ic.EchoHttpServer.GetEchoInstance().Group(ic.EchoHttpServer.GetBasePath())
	httpController := authHttp.NewController(useCase)
	authHttp.NewRouter(httpController).Register(httpRouterGp)

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
