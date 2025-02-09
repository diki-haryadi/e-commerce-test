package productFixture

import (
	"context"
	"math"
	"net"
	"time"

	productV1 "github.com/diki-haryadi/protobuf-template/go-micro-template/article/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	sampleExtServiceUseCase "github.com/diki-haryadi/go-micro-template/external/sample_ext_service/usecase"
	productGrpc "github.com/diki-haryadi/go-micro-template/internal/warehouse/delivery/grpc"
	productHttp "github.com/diki-haryadi/go-micro-template/internal/warehouse/delivery/http"
	productKafkaProducer "github.com/diki-haryadi/go-micro-template/internal/warehouse/delivery/kafka/producer"
	productRepo "github.com/diki-haryadi/go-micro-template/internal/warehouse/repository"
	productUseCase "github.com/diki-haryadi/go-micro-template/internal/warehouse/usecase"
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
	productGrpcClient productV1.productServiceClient
}

func NewIntegrationTestFixture() (*IntegrationTestFixture, error) {
	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	ic, infraDown, err := iContainer.NewIC(ctx)
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

	// grpc
	grpcController := productGrpc.NewController(useCase)
	productV1.RegisterproductServiceServer(ic.GrpcServer.GetCurrentGrpcServer(), grpcController)

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

	productGrpcClient := productV1.NewProductServiceClient(grpcClientConn)

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
		productGrpcClient: productGrpcClient,
	}, nil
}
