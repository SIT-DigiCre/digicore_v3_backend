package google

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type ResponseOAuthURL struct {
	URL string `json:"url"`
}

// Get OAuth request url
// @Router /google/oauth/url [get]
// @Param register query string false "type"
// @Success 302 "redirect to oauth page"
// @Header 302 {string}  Location
func (c Context) OAuthURL(e echo.Context) error {
	register := e.QueryParam("register")
	if register == "true" {
		redirectURL := oauth2.SetAuthURLParam("redirect_uri", env.BackRootURL+"/google/oauth/callback/register")
		authURL := c.Config.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce, redirectURL)
		return e.Redirect(http.StatusFound, authURL)
	}
	redirectURL := oauth2.SetAuthURLParam("redirect_uri", env.BackRootURL+"/google/oauth/callback/login")
	authURL := c.Config.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce, redirectURL)
	return e.Redirect(http.StatusFound, authURL)
}
