package action

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func Run(
	ctx context.Context,
	before func(ctx context.Context) error,
	rpcRegister func(grpc.ServiceRegistrar) error,
	rpcGatewayRegister func(*runtime.ServeMux, string, []grpc.DialOption) error,
	watch func(ctx context.Context) error,
) error {
	if before != nil {
		if err := before(ctx); err != nil {
			logger.Sugar().Errorw("Run", "Before", err)
			return err
		}
	}

	// https://pkg.go.dev/syscall#SIGINT
	ctx, stop := signal.NotifyContext(
		ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGABRT,
		syscall.SIGILL,
		syscall.SIGBUS,
		syscall.SIGFPE,
		syscall.SIGPIPE,
		syscall.SIGQUIT,
		syscall.SIGSEGV,
		syscall.SIGTERM,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGTSTP,
	)

	if watch != nil {
		if err := watch(ctx); err != nil {
			logger.Sugar().Errorw("Run", "Watch", err)
			return err
		}
	}

	go func() {
		if err := grpc2.RunGRPC(rpcRegister); err != nil {
			logger.Sugar().Errorw("Run", "GRPCRegister", err)
		}
	}()
	go func() {
		if err := grpc2.RunGRPCGateWay(rpcGatewayRegister); err != nil {
			logger.Sugar().Errorw("Run", "GRPCGatewayRegister", err)
		}
	}()

	<-ctx.Done()
	if ctx.Err() != nil {
		logger.Sugar().Errorw("Run", "Exit", ctx.Err())
	}
	stop()

	if err := grpc2.HShutdown(); err != nil {
		logger.Sugar().Warnw("Run", "GRPCGatewayShutdown", err)
	}
	grpc2.GShutdown()

	return nil
}
