package discord

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
)

// OAuth callback destination
// @Router /discord/oauth/callback [get]
// @Param code query string true "oauth code"
// @Success 302 "send oauth code to frontend"
// @Header 302 {string}  Location "/user/discord?code={}"
func (c Context) OAuthCallback(e echo.Context) error {
	code := e.QueryParam("code")
	return e.Redirect(http.StatusFound, env.FrontendRootURL+"/user/discord?code="+code)
}
