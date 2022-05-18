package admin

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

func Middleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := user.GetUserId(&c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest)
			}
			err = db.QueryRow("SELECT id FROM GroupUser LEFT JOIN User ON GroupUser.user_id = User.id WHERE User.id = ? AND group_id = ?", id, env.AdminGroupID).Scan()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			return next(c)
		}
	}
}
