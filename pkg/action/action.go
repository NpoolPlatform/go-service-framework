package action

import (
	"context"
	"os"
	"os/signal"
	goruntime "runtime"
	"syscall"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// nolint
func Run(
	ctx context.Context,
	init func(ctx context.Context) error,
	rpcRegister func(grpc.ServiceRegistrar) error,
	rpcGatewayRegister func(*runtime.ServeMux, string, []grpc.DialOption) error,
	watch func(ctx context.Context, cancel context.CancelFunc) error,
	rpcSecureRegister *func(grpc.ServiceRegistrar) error,
) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	signal.Ignore(syscall.SIGPIPE)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	runRPC := func(rpcRegister func(grpc.ServiceRegistrar) error, secure bool) {
		defer func() {
			if err := recover(); err != nil {
				logger.Sugar().Errorw(
					"Watch",
					"State", "Panic",
					"Error", err,
				)
				cancel()
			}
		}()
		if err := grpc2.RunGRPC(rpcRegister, func(p interface{}) error {
			const defaultStackSize = 8192
			var buf [defaultStackSize]byte
			n := goruntime.Stack(buf[:], false)
			logger.Sugar().Errorw(
				"Watch",
				"State", "Panic",
				"P", p,
				"Stack", string(buf[:n]),
			)
			cancel()
			name := config.GetStringValueWithNameSpace("", config.KeyHostname)
			return wlog.Errorf("Panic (%v): %v", name, p)
		}, secure); err != nil {
			logger.Sugar().Errorw("Run", "GRPCRegister", err)
		}
	}

	go runRPC(rpcRegister, false)
	if rpcSecureRegister != nil {
		go runRPC(*rpcSecureRegister, true)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Sugar().Errorw(
					"Watch",
					"State", "Panic",
					"Error", err,
				)
				cancel()
			}
		}()
		if err := grpc2.RunGRPCGateWay(rpcGatewayRegister); err != nil {
			logger.Sugar().Errorw("Run", "GRPCGatewayRegister", err)
		}
	}()

	go func() {
		defer cancel()
		for {
			sig := <-sigs
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
			case syscall.SIGQUIT:
				fallthrough //nolint
			case syscall.SIGSEGV:
				fallthrough //nolint
			case syscall.SIGTERM:
				logger.Sugar().Warnw("Run", "Exit", sig)
				return
			case syscall.SIGPIPE:
				logger.Sugar().Warnw("Run", "Exception", sig)
			}
		}
	}()

	if init != nil {
		if err := init(ctx); err != nil {
			logger.Sugar().Errorw("Run", "Before", err)
			return wlog.WrapError(err)
		}
	}

	if watch != nil {
		if err := watch(ctx, cancel); err != nil {
			logger.Sugar().Errorw("Run", "Watch", err)
			return wlog.WrapError(err)
		}
	}

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

func Watch(
	ctx context.Context,
	cancel context.CancelFunc,
	w func(ctx context.Context),
	p func(ctx context.Context),
) {
	defer func() {
		if err := recover(); err != nil {
			const defaultStackSize = 8192
			var buf [defaultStackSize]byte
			n := goruntime.Stack(buf[:], false)
			logger.Sugar().Errorw(
				"Watch",
				"State", "Panic",
				"Error", err,
				"Stack", string(buf[:n]),
			)
			p(ctx)
			cancel()
		}
	}()
	w(ctx)
}
