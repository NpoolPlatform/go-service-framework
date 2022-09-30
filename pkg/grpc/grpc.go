package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	// ErrServiceIDEmpty ..
	ErrServiceIDEmpty = errors.New("service id empty")
	// ErrServiceIDInvalid ..
	ErrServiceIDInvalid = errors.New("service id invalid uuid")
)

const (
	GRPCTAG = "GRPCTAG"
	HTTPTAG = "HTTPTAG"
)

var (
	grpcServer       *grpc.Server
	httpServer       *http.Server
	jaegerTp         *trace.TracerProvider
	registerDuration = 10 * time.Second
)

func GShutdown() {
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
}

func HShutdown() error {
	if httpServer != nil {
		return httpServer.Shutdown(context.Background())
	}
	return nil
}

func TShutdown() error {
	if jaegerTp != nil {
		return jaegerTp.Shutdown(context.Background())
	}
	return nil
}

func registerConsul(healthCheck bool, id, name, tag string, port int) {
	hp := 0
	if healthCheck {
		hp = port
	}

	for range time.NewTicker(registerDuration).C {
		err := consul.RegisterService(healthCheck, consul.RegisterInput{
			ID:          id,
			Name:        name,
			Tags:        []string{tag},
			Port:        port,
			HealthzPort: hp,
		})
		if err != nil {
			logger.Sugar().Errorf("fail to register consul service: %v", err)
		}
	}
}

func RunGRPC(serviceRegister func(srv grpc.ServiceRegistrar) error) error {
	if serviceRegister == nil {
		return xerrors.Errorf("service register must be set")
	}

	gport := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)
	name := config.GetStringValueWithNameSpace("", config.KeyHostname)
	prometheusPort := config.GetIntValueWithNameSpace("", config.KeyPrometheusPort)

	var err error
	// peek collect service endpoint

	// init jaeger provider
	jaegerTp, err = jaegerTracerProvider(
		// here use sider car
		// jaeger-agent.kube-system.svc.cluster.local
		"127.0.0.1",
		"6831",
		config.GetStringValueWithNameSpace("", config.KeyENV),
		config.GetStringValueWithNameSpace("", config.KeyHostname),
		config.GetStringValueWithNameSpace("", config.KeyServiceID),
	)
	if err != nil {
		return xerrors.Errorf("fail to init tracer %v", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", gport))
	if err != nil {
		return xerrors.Errorf("fail to listen tcp at %v: %v", gport, err)
	}

	grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			otelgrpc.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			otelgrpc.UnaryServerInterceptor(),
		)),
	)

	sid := config.GetStringValueWithNameSpace("", config.KeyServiceID)
	if sid == "" {
		return ErrServiceIDEmpty
	}
	if _, err := uuid.Parse(sid); err != nil {
		return ErrServiceIDInvalid
	}

	go registerConsul(
		false,
		fmt.Sprintf("%s-%s", GRPCTAG, sid),
		name,
		GRPCTAG,
		gport,
	)

	err = serviceRegister(grpcServer)
	if err != nil {
		return xerrors.Errorf("fail to register services: %v", err)
	}

	reflection.Register(grpcServer)

	// prometheus metrics endpoints
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%v", prometheusPort), nil) //nolint
	}()

	return grpcServer.Serve(l)
}

func RunGRPCGateWay(serviceRegister func(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error) error {
	if serviceRegister == nil {
		return xerrors.Errorf("service register must be set")
	}

	gport := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)
	hport := config.GetIntValueWithNameSpace("", config.KeyHTTPPort)
	name := config.GetStringValueWithNameSpace("", config.KeyHostname)
	mux := runtime.NewServeMux()
	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%v", hport),
		Handler: mux,
	}

	// consul health check
	if err := mux.HandlePath(http.MethodGet, "/healthz", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("PONG")) // nolint
	}); err != nil {
		return xerrors.Errorf("fail to healthz check: %v", err)
	}

	sid := config.GetStringValueWithNameSpace("", config.KeyServiceID)
	if sid == "" {
		return ErrServiceIDEmpty
	}
	if _, err := uuid.Parse(sid); err != nil {
		return ErrServiceIDInvalid
	}

	go registerConsul(
		true,
		fmt.Sprintf("%s-%s", HTTPTAG, sid),
		name,
		HTTPTAG,
		hport,
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := serviceRegister(mux, fmt.Sprintf(":%v", gport), opts)
	if err != nil {
		return xerrors.Errorf("fail to register services: %v", err)
	}

	return httpServer.ListenAndServe()
}

// GetGRPCConn get grpc client conn
func GetGRPCConn(service string, tags ...string) (*grpc.ClientConn, error) {
	if service == "" {
		return nil, fmt.Errorf("service is empty")
	}

	svc, err := config.PeekService(service, tags...)
	if err != nil {
		return nil, err
	}

	targets := strings.Split(
		net.JoinHostPort(svc.Address, fmt.Sprintf("%d", svc.Port)), ",")

	for _, target := range targets {
		_ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

		conn, err := grpc.DialContext(_ctx, target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock())

		cancel()

		if err != nil {
			logger.Sugar().Errorf("fail to dial grpc %v: %v", target, err)
			continue
		}

		connState := conn.GetState()
		if connState != connectivity.Idle && connState != connectivity.Ready {
			logger.Sugar().Warnf("conn not available %v: %v", target, connState)
			continue
		}

		return conn, nil
	}

	return nil, fmt.Errorf("valid conn not found")
}
