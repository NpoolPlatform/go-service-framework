package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"
	http2 "github.com/NpoolPlatform/go-service-framework/pkg/http"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cli "github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Action: func(c *cli.Context) error {
		go func() {
			err := grpc2.Run(rpcRegister)
			if err != nil {
				logger.Sugar().Errorf("fail to run grpc server: %v", err)
			}
		}()
		return http2.Run(registerRoute)
	},
}

func registerRoute(router *chi.Mux) error {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http2.Response(w, []byte("hello chi"), 0, "") //nolint
	})
	return nil
}

func rpcRegister(server *grpc.Server) error {
	return nil
}
