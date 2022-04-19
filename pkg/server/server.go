package server

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/discord"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateEchoServer(db *sql.DB) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	addRouting(e, db)
	return e
}

func addRouting(e *echo.Echo, db *sql.DB) {
	google, _ := google.CreateContext(db)
	e.GET("/google/oauth/url", google.OAuthURL)
	e.GET("/google/oauth/callback/login", google.OAuthCallbackLogin)
	e.GET("/google/oauth/callback/register", google.OAuthCallbackRegister)

	discord, _ := discord.CreateContext()
	e.GET("/discord/oauth/url", discord.OAuthURL)
	e.GET("/discord/oauth/callback", discord.OAuthCallback)

	config := middleware.JWTConfig{
		SigningKey: []byte(env.JWTSecret),
		ErrorHandler: func(error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
		},
	}

	r := e.Group("/user")
	r.Use(middleware.JWTWithConfig(config))
	user, _ := user.CreateContext(db)
	r.PUT("/my", user.UpdateMyProfile)
	r.GET("/my", user.GetMyProfile)
	r.PUT("/my/discord", user.UpdateDiscordId)
	r.PUT("/my/private", user.UpdateMyPrivateProfile)
	r.GET("/my/private", user.GetMyPrivateProfile)
	r.PUT("/my/payment", user.UpdateMyPayment)
	r.GET("/my/payment", user.GetMyPayment)
	r.GET("/my/payment/history", user.GetMyPaymentHistory)
}

func CreateDbConnection(address string) (*sql.DB, error) {
	db, err := sql.Open("mysql", address)
	return db, err
}
