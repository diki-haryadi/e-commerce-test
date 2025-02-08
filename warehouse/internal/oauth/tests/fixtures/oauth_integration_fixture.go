package oauthFixture

import (
	"context"
	"math"
	"net"
	"time"

	articleV1 "github.com/diki-haryadi/protobuf-template/go-micro-template/article/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	oauthGrpc "github.com/diki-haryadi/go-micro-template/internal/oauth/delivery/grpc"
	oauthHttp "github.com/diki-haryadi/go-micro-template/internal/oauth/delivery/http"
	oauthKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/oauth/delivery/kafka/producer"
	oauthRepo "github.com/diki-haryadi/go-micro-template/internal/oauth/repository"
	oauthUseCase "github.com/diki-haryadi/go-micro-template/internal/oauth/usecase"
	externalBridge "github.com/diki-haryadi/ztools/external_bridge"
	iContainer "github.com/diki-haryadi/ztools/infra_container"
	"github.com/diki-haryadi/ztools/logger"
)

const BUFSIZE = 1024 * 1024

type IntegrationTestFixture struct {
	TearDown          func()
	Ctx               context.Context
	Cancel            context.CancelFunc
	InfraContainer    *iContainer.IContainer
	ArticleGrpcClient articleV1.ArticleServiceClient
}

func NewIntegrationTestFixture() (*IntegrationTestFixture, error) {
	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	container := iContainer.IContainer{}
	ic, infraDown, err := container.IContext(ctx).
		ICDown().ICPostgres().ICGrpc().ICEcho().
		ICKafka().NewIC()
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
	kafkaProducer := oauthKafkaProducer.NewProducer(ic.KafkaWriter)
	repository := oauthRepo.NewRepository(ic.Postgres)
	useCase := oauthUseCase.NewUseCase(repository, seServiceUseCase, kafkaProducer)

	// http
	ic.EchoHttpServer.SetupDefaultMiddlewares()
	httpRouterGp := ic.EchoHttpServer.GetEchoInstance().Group(ic.EchoHttpServer.GetBasePath())
	httpController := oauthHttp.NewController(useCase)
	oauthHttp.NewRouter(httpController).Register(httpRouterGp)

	// grpc
	grpcController := oauthGrpc.NewController(useCase)
	articleV1.RegisterArticleServiceServer(ic.GrpcServer.GetCurrentGrpcServer(), grpcController)

	lis := bufconn.Listen(BUFSIZE)
	go func() {
		if err := ic.GrpcServer.GetCurrentGrpcServer().Serve(lis); err != nil {
			logger.Zap.Sugar().Fatalf("Server exited with error: %v", err)
		}
	}()

	grpcClientConn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		cancel()
		return nil, err
	}

	articleGrpcClient := articleV1.NewArticleServiceClient(grpcClientConn)

	return &IntegrationTestFixture{
		TearDown: func() {
			cancel()
			infraDown()
			_ = grpcClientConn.Close()
			extBridgeDown()
		},
		InfraContainer:    ic,
		Ctx:               ctx,
		Cancel:            cancel,
		ArticleGrpcClient: articleGrpcClient,
	}, nil
}
