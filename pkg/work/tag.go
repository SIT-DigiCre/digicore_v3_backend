package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseUpdateDiscordId struct {
	Error string `json:"error"`
}

// Create tag
// @Accept json
// @Security Authorization
// @Router /work/tag [post]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) CreateTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Update tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) UpdateTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Get tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [get]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) GetTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Get tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [delete]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) DeleteTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}
