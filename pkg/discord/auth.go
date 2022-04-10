package discord

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
)

// Get OAuth request url
// @Accept json
// @Router /google/oauth/url [get]
// @Success 200 {object} ResponseOAuthURL
func (c Context) OAuthURL(e echo.Context) error {
	return e.Redirect(http.StatusFound, fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify", env.DiscordClientID, env.DiscordRedirectURL))
}
