package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) WorkList(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) TagList(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}
