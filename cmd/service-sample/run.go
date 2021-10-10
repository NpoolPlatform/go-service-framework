package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/NpoolPlatform/go-service-framework/pkg/server"

	cli "github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Action: func(c *cli.Context) error {
		return server.Run(registerRoute)
	},
}

func registerRoute(router *chi.Mux) error {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello chi")) //nolint
	})
	return nil
}
