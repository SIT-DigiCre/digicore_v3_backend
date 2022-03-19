package user

import (
	"database/sql"
	"fmt"

	echo_session "github.com/ipfans/echo-session"
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
	session := echo_session.Default(*e)
	_, ok := session.Get("login").(bool)
	if !ok {
		return "", fmt.Errorf("aaa")
	}
	id, ok := session.Get("id").(string)
	if !ok {
		return "", fmt.Errorf("bbb")
	}
	return id, nil
}
