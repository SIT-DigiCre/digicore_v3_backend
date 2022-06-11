package admin

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type Context struct {
	DB *sql.DB
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}

func Middleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := user.GetUserId(&c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest)
			}
			count := 0
			err = db.QueryRow("SELECT count(*) FROM groups_users LEFT JOIN users ON groups_users.user_id = users.id WHERE users.id = UUID_TO_BIN(?) AND group_id = UUID_TO_BIN(?)", id, env.AdminGroupID).Scan(&count)
			if err != nil || count == 0 {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			return next(c)
		}
	}
}
