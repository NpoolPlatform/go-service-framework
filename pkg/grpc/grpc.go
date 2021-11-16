package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	GRPCTAG = "GRPCTAG"
	HTTPTAG = "HTTPTAG"
)

var (
	target2Conn sync.Map
	grpcServer  *grpc.Server
	httpServer  *http.Server
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

func RunGRPC(serviceRegister func(srv grpc.ServiceRegistrar) error) error {
	if serviceRegister == nil {
		return xerrors.Errorf("service register must be set")
	}

	gport := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)
	name := config.GetStringValueWithNameSpace("", config.KeyHostname)
	prometheusPort := config.GetIntValueWithNameSpace("", config.KeyPrometheusPort)

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", gport))
	if err != nil {
		return xerrors.Errorf("fail to listen tcp at %v: %v", gport, err)
	}

	grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)

	err = consul.RegisterService(false, consul.RegisterInput{
		ID:   uuid.New(),
		Name: name,
		Tags: []string{GRPCTAG},
		Port: gport,
	})
	if err != nil {
		return xerrors.Errorf("fail to register consul service: %v", err)
	}

	err = serviceRegister(grpcServer)
	if err != nil {
		return xerrors.Errorf("fail to register services: %v", err)
	}

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

	err := consul.RegisterService(true, consul.RegisterInput{
		ID:          uuid.New(),
		Name:        name,
		Tags:        []string{HTTPTAG},
		Port:        hport,
		HealthzPort: hport,
	})
	if err != nil {
		return xerrors.Errorf("fail to register consul service: %v", err)
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = serviceRegister(mux, fmt.Sprintf(":%v", gport), opts)
	if err != nil {
		return xerrors.Errorf("fail to register services: %v", err)
	}

	return httpServer.ListenAndServe()
}

// GetGRPCConn get grpc client conn
func GetGRPCConn(address string) (*grpc.ClientConn, error) {
	if address == "" {
		return nil, fmt.Errorf("address is empty")
	}

	targets := strings.Split(address, ",")

	for _, target := range targets {
		v, ok := target2Conn.Load(target)
		if !ok {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(),
				grpc.WithBlock(),
			)
			if err != nil {
				continue
			}
			target2Conn.Store(target, conn)
			return conn, nil
		}

		var conn *grpc.ClientConn
		if _conn, ok := v.(*grpc.ClientConn); ok {
			conn = _conn
		}
		if conn == nil {
			continue
		}

		connState := conn.GetState()
		if connState != connectivity.Idle && connState != connectivity.Ready {
			continue
		}
		return conn, nil
	}
	return nil, fmt.Errorf("valid conn not found")
}
