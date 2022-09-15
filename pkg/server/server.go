package server

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/discord"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/event"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/storage"
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

	group_group := e.Group("/group")
	group_group.Use(middleware.JWTWithConfig(config))
	group, _ := group.CreateContext(db)
	group_group.GET("", group.GroupList)
	group_group.GET("/:id", group.Detail)
	group_group.POST("/:id", group.Join)
	group_group.DELETE("/:id", group.Leave)

	user_group := e.Group("/user")
	user_group.Use(middleware.JWTWithConfig(config))
	user, _ := user.CreateContext(db)
	user_group.GET("", user.GetList)
	user_group.PUT("/my", user.UpdateMyProfile)
	user_group.GET("/my", user.GetMyProfile)
	user_group.PUT("/my/introduction", user.UpdateMySelfIntroduction)
	user_group.GET("/my/introduction", user.GetMySelfIntroduction)
	user_group.PUT("/my/discord", user.UpdateDiscordId)
	user_group.PUT("/my/private", user.UpdateMyPrivateProfile)
	user_group.GET("/my/private", user.GetMyPrivateProfile)
	user_group.PUT("/my/payment", user.UpdateMyPayment)
	user_group.GET("/my/payment", user.GetMyPayment)
	user_group.GET("/my/payment/history", user.GetMyPaymentHistory)
	user_group.GET("/:id", user.GetProfile)
	user_group.GET("/:id/introduction", user.GetSelfIntroduction)

	event_group := e.Group("/event")
	event_group.Use(middleware.JWTWithConfig(config))
	event, _ := event.CreateContext(db)
	event_group.GET("", event.GetEventsList)
	event_group.GET("/:id", event.GetEventDetail)
	event_group.GET("/:event_id/:id", event.ReservationInfo)
	event_group.POST("/:event_id/:id", event.Reservation)
	event_group.DELETE("/:event_id/:id", event.CancelReservation)

	s := e.Group("/storage")
	s.Use(middleware.JWTWithConfig(config))
	storage, _ := storage.CreateContext(db)
	s.POST("", storage.UploadUserfile)
	s.GET("/:fileId", storage.GetUserFileUrl)

	env_group := e.Group("/env")
	env_group.Use(middleware.JWTWithConfig(config))
	env, _ := env.CreateContext()
	env_group.GET("/join", env.GetJoinURL)
}

func CreateDbConnection(address string) (*sql.DB, error) {
	db, err := sql.Open("mysql", address)
	return db, err
}
