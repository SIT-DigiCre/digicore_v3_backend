package env

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseJoinURL struct {
	SlackURL   string `json:"slack_url"`
	DiscordURL string `json:"discord_url"`
	Error      string `json:"error"`
}

// Get env url
// @Router /env/join [get]
// @Security Authorization
// @Success 200 {object} ResponseJoinURL
func (c Context) GetJoinURL(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseJoinURL{SlackURL: SlackURL, DiscordURL: DiscordURL})
}
