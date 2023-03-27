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
	init func(ctx context.Context) error,
	rpcRegister func(grpc.ServiceRegistrar) error,
	rpcGatewayRegister func(*runtime.ServeMux, string, []grpc.DialOption) error,
	watch func(ctx context.Context) error,
) error {
	if init != nil {
		if err := init(ctx); err != nil {
			logger.Sugar().Errorw("Run", "Before", err)
			return err
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	go func() {
		defer cancel()
		for {
			sig := <-sigs
			logger.Sugar().Infow("Run", "Signal", sig)
			switch sig {
			case syscall.SIGKILL:
				fallthrough //nolint
			case syscall.SIGABRT:
				fallthrough //nolint
			case syscall.SIGBUS:
				fallthrough //nolint
			case syscall.SIGFPE:
				fallthrough //nolint
			case syscall.SIGILL:
				fallthrough //nolint
			case syscall.SIGINT:
				fallthrough //nolint
			case syscall.SIGPIPE:
				fallthrough //nolint
			case syscall.SIGQUIT:
				fallthrough //nolint
			case syscall.SIGSEGV:
				fallthrough //nolint
			case syscall.SIGTERM:
				logger.Sugar().Infow("Run", "Exit", sig)
				return
			}
		}
	}()

	<-ctx.Done()
	if ctx.Err() != nil {
		logger.Sugar().Errorw("Run", "Exit", ctx.Err())
	}

	if err := grpc2.HShutdown(); err != nil {
		logger.Sugar().Warnw("Run", "GRPCGatewayShutdown", err)
	}
	grpc2.GShutdown()

	return nil
}
