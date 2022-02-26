package google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseOAuthURL struct {
	URL string `json:"url"`
}

// Get OAuth request url
// @Accept json
// @Router /google/oauth/url [get]
// @Success 200 {object} ResponseOAuthURL
func (c Context) OAuthURL(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseOAuthURL{})
}
