package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google"
	echo_session "github.com/ipfans/echo-session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateEchoServer(store echo_session.RedisStore) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(echo_session.Sessions("SESSION", store))
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

func CreateSessionStoreConnection(address string, password string) (echo_session.RedisStore, error) {
	store, err := echo_session.NewRedisStore(32, "tcp", address, password, make([]byte, 32))
	if err != nil {
		return nil, err
	}
	store.Options(echo_session.Options{Path: "/", MaxAge: 86400 * 7})

	return store, nil
}
