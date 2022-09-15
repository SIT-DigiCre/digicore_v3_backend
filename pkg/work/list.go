package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Get work list
// @Accept json
// @Security Authorization
// @Router /work/work [get]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) WorkList(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Get tag list
// @Accept json
// @Security Authorization
// @Router /work/tag [get]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) TagList(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}
