package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseAuth struct {
	Error string `json:"error"`
}

// Admin user auth
// @Router /admin [get]
// @Security Authorization
// @Success 200 {object} ResponseAuth
func (c Context) Auth(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseAuth{})
}
