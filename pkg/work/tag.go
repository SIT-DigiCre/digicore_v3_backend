package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestCreateTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseCreateTag struct {
	Success bool `json:"success"`
}

type RequestUpdatgeTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseUpdatgeTag struct {
	Success bool `json:"success"`
}

type ResponseGetTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ResponseDeleteTag struct {
	Success bool `json:"success"`
}

// Create tag
// @Accept json
// @Security Authorization
// @Router /work/tag [post]
// @Param RequestCreateTag body RequestCreateTag true "update work"
// @Success 200 {object} ResponseCreateTag
func (c Context) CreateTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseCreateTag{})
}

// Update tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [put]
// @Param id path string true "tag id"
// @Param RequestUpdatgeTag body RequestUpdatgeTag true "update work"
// @Success 200 {object} ResponseUpdatgeTag
func (c Context) UpdateTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseUpdatgeTag{})
}

// Get tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [get]
// @Param id path string true "tag id"
// @Success 200 {object} ResponseGetTag
func (c Context) GetTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseGetTag{})
}

// Get tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [delete]
// @Param id path string true "tag id"
// @Success 200 {object} ResponseDeleteTag
func (c Context) DeleteTag(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseDeleteTag{})
}
