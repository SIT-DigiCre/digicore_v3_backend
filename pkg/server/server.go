package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateEchoServer() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	addRouting(e)
	return e
}

func addRouting(e *echo.Echo) {
	google, _ := google.CreateContext()
	e.GET("/google/oauth/url", google.OAuthURL)
	e.GET("/google/oauth/callback", google.OAuthCallback)
}
