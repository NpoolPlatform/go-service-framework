package http

import (
	"encoding/json"
	"fmt"
	http2 "net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

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

	r.Get("/healthz", func(w http2.ResponseWriter, r *http2.Request) {
		Response(w, []byte("pong"), 0, "") //nolint
	})

	err = consul.RegisterService(true, consul.RegisterInput{
		ID:          uuid.New().String(),
		Name:        config.GetStringValueWithNameSpace("", config.KeyHostname),
		Tags:        nil,
		Port:        config.GetIntValueWithNameSpace("", config.KeyHTTPPort),
		HealthzPort: config.GetIntValueWithNameSpace("", config.KeyHTTPPort),
	})
	if err != nil {
		return xerrors.Errorf("fail to register service: %v", err)
	}

	listen := fmt.Sprintf(":%v", config.GetIntValueWithNameSpace("", config.KeyHTTPPort))
	logger.Sugar().Infof("Start HTTP server: %v", listen)

	return http2.ListenAndServe(listen, r)
}

func Response(w http2.ResponseWriter, body interface{}, code int, msg string) error { //nolint
	type input struct {
		Body interface{} `json:"body"`
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
	}

	in := input{
		Body: body,
		Code: code,
		Msg:  msg,
	}

	jsonStr, err := json.Marshal(&in)
	if err != nil {
		return xerrors.Errorf("fail to marshal body: %v", err)
	}

	_, err = w.Write(jsonStr)
	return err
}
