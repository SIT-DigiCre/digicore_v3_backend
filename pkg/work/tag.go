package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseUpdateDiscordId struct {
	Error string `json:"error"`
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) CreateTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) UpdateTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) GetTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) DeleteTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}
