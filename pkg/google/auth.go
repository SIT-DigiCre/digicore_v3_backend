package google

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type ResponseOAuthURL struct {
	URL string `json:"url"`
}

// Get OAuth request url
// @Accept json
// @Router /google/oauth/url [get]
// @Success 200 {object} ResponseOAuthURL
func (c Context) OAuthURL(e echo.Context) error {
	authURL := c.Config.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return e.JSON(http.StatusOK, ResponseOAuthURL{URL: authURL})
}
