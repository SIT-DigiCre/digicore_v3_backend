package google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseOAuthCallback struct {
}

func (c Context) OAuthCallback(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseOAuthCallback{})
}
