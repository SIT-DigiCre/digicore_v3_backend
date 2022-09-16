package work

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Work struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Authers []Auther `json:"authers"`
	Tags    []Tag    `json:"tags"`
}

type ResponseGetWorkList struct {
	Works []Work `json:"works"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResponseGetTagList struct {
	Tags []Tag `json:"tags"`
}

// Get work list
// @Accept json
// @Security Authorization
// @Router /work/work [get]
// @Success 200 {object} ResponseGetWorkList
func (c Context) WorkList(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseGetWorkList{})
}

// Get tag list
// @Accept json
// @Security Authorization
// @Router /work/tag [get]
// @Success 200 {object} ResponseGetTagList
func (c Context) TagList(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseGetTagList{})
}
