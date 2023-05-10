package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/authenticator"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func CreateEchoServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	authenticator_middleware, err := authenticator.CreateAuthenticator()
	if env.Auth == "disable" {
		authenticator_middleware, err = authenticator.CreateDebugAuthenticator()
	}
	if err != nil {
		logrus.Fatal("Failed to create validation middleware: %w", err)
	}
	e.Use(authenticator_middleware...)

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	server := NewServer()
	api.RegisterHandlers(e, server)
	return e
}
