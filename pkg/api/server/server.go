package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/validator"
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

	validater_middleware, err := validator.CreateValidator()
	if err != nil {
		logrus.Fatal("Failed to create validation middleware: %w", err)
	}
	e.Use(validater_middleware...)

	server := NewServer()
	api.RegisterHandlers(e, server)
	return e
}
