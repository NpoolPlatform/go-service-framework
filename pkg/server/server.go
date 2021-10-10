package server

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(routeRegister func(router *chi.Mux) error) error {
	if routeRegister == nil {
		return xerrors.Errorf("ROUTE REGISTER CALLBACK IS MUST")
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	err := routeRegister(r)
	if err != nil {
		return xerrors.Errorf("fail to register route: %v", err)
	}

	listen := fmt.Sprintf(":%v", config.GetIntValueWithNameSpace("", config.KeyHTTPPort))
	return http.ListenAndServe(listen, r)
}
