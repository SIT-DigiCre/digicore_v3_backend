package discord

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
)

// OAuth callback destination
// @Accept json
// @Router /discord/oauth/callback [get]
// @Param code query string true "auth token"
// @Success 302
func (c Context) OAuthCallback(e echo.Context) error {
	code := e.QueryParam("code")
	return e.Redirect(http.StatusFound, env.FrontRootURL+"/user/discord?code="+code)
}
