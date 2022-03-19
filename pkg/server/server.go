package server

import (
	"database/sql"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/google"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	_ "github.com/go-sql-driver/mysql"
	echo_session "github.com/ipfans/echo-session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateEchoServer(store echo_session.RedisStore, db *sql.DB) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(echo_session.Sessions("session", store))
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	addRouting(e, db)
	return e
}

func addRouting(e *echo.Echo, db *sql.DB) {
	google, _ := google.CreateContext(db)
	e.GET("/google/oauth/url", google.OAuthURL)
	e.GET("/google/oauth/callback", google.OAuthCallback)

	user, _ := user.CreateContext(db)
	e.POST("/user/my/private", user.SetMyPrivateProfile)
	e.GET("/user/my/private", user.GetMyPrivateProfile)
}

func CreateDbConnection(address string) (*sql.DB, error) {
	db, err := sql.Open("mysql", address)
	return db, err
}

func CreateSessionStoreConnection(address string, password string) (echo_session.RedisStore, error) {
	store, err := echo_session.NewRedisStore(32, "tcp", address, password, make([]byte, 32))
	if err != nil {
		return nil, err
	}
	store.Options(echo_session.Options{Path: "/", MaxAge: 86400 * 7})

	return store, nil
}
