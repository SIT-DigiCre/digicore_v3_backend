package google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseOAuthURL struct {
	URL string `json:"url"`
}

func (c Context) OAuthURL(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseOAuthURL{})
}
