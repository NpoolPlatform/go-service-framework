package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	http2 "github.com/NpoolPlatform/go-service-framework/pkg/http"

	cli "github.com/urfave/cli/v2"
)

var runCmd = &cli.Command{
	Name:    "run",
	Aliases: []string{"s"},
	Usage:   "Run the daemon",
	Action: func(c *cli.Context) error {
		return http2.Run(registerRoute)
	},
}

func registerRoute(router *chi.Mux) error {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http2.Response(w, []byte("hello chi"), 0, "") //nolint
	})
	return nil
}
