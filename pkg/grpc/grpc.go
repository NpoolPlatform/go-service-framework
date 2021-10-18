package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var target2Conn sync.Map

func Run(serviceRegister func(srv grpc.ServiceRegistrar) error) error {
	if serviceRegister == nil {
		return xerrors.Errorf("service register must be set")
	}

	port := config.GetIntValueWithNameSpace("", config.KeyGRPCPort)
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", config.GetIntValueWithNameSpace("", "grpc_port")))
	if err != nil {
		return xerrors.Errorf("fail to listen tcp at %v: %v", port, err)
	}

	srv := grpc.NewServer()
	err = serviceRegister(srv)
	if err != nil {
		return xerrors.Errorf("fail to register services: %v", err)
	}

	return srv.Serve(l)
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
