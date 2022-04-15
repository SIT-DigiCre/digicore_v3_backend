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
// @Router /google/oauth/url [get]
// @Success 302 "redirect to oauth page"
// @Header 302 {string}  Location
func (c Context) OAuthURL(e echo.Context) error {
	authURL := c.Config.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return e.Redirect(http.StatusFound, authURL)
}
