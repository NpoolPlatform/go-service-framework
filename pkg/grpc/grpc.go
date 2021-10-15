package grpc

import (
	"fmt"
	"net"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"

	"google.golang.org/grpc"
)

func Run(serviceRegister func(srv *grpc.Server) error) error {
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
