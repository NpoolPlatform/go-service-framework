package main

import (
	"google.golang.org/grpc"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	cli "github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "Run the daemon",
	After: func(c *cli.Context) error {
		// close db, http or grpc server graceful shutdown
		if err := grpc2.HShutdown(); err != nil {
			return err
		}
		grpc2.GShutdown()
		return logger.Sync()
	},
	Action: func(c *cli.Context) error {
		go func() {
			err := grpc2.RunGRPC(rpcRegister, func(p interface{}) error {
				return nil
			})
			if err != nil {
				logger.Sugar().Errorf("fail to run grpc server: %v", err)
			}
		}()

		return grpc2.RunGRPCGateWay(rpcGatewayRegister)
	},
}

func rpcRegister(server grpc.ServiceRegistrar) error {
	return nil
}

func rpcGatewayRegister(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}
