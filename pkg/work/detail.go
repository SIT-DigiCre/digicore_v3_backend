package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResGetWork struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Authers     []string `json:"authers"`
}

// Create work
// @Accept json
// @Security Authorization
// @Router /work/work [post]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) CreateWork(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Update work
// @Accept json
// @Security Authorization
// @Router /work/work/{id} [post]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) UpdateWork(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Get work
// @Accept json
// @Security Authorization
// @Router /work/work/{id} [get]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) GetWork(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

// Delete work
// @Accept json
// @Security Authorization
// @Router /work/work/{id} [delete]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) DeleteWork(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}
