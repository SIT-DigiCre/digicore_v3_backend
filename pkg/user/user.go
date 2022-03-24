package user

import (
	"database/sql"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Context struct {
	DB *sql.DB
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}

func GetUserId(e *echo.Context) (string, error) {
	user := (*e).Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["uuid"].(string)
	return id, nil
}
