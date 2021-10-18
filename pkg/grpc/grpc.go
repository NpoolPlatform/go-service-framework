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
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
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

	port := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return xerrors.Errorf("fail to listen tcp at %v: %v", port, err)
	}

	grpcServer = grpc.NewServer()
	err = serviceRegister(grpcServer)
	if err != nil {
		return xerrors.Errorf("fail to register services: %v", err)
	}

	return grpcServer.Serve(l)
}

func RunGRPCGateWay(serviceRegister func(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error) error {
	if serviceRegister == nil {
		return xerrors.Errorf("service register must be set")
	}

	gport := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)
	hport := config.GetIntValueWithNameSpace("", config.KeyHTTPPort)

	mux := runtime.NewServeMux()
	httpServer = &http.Server{
		Addr:    fmt.Sprintf("%v", hport),
		Handler: mux,
	}
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := serviceRegister(mux, fmt.Sprintf("%v", gport), opts)
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
