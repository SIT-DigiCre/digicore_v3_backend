package work

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestCreateTag = UpdatgeTag

type ResponseCreateTag struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type RequestUpdatgeTag = UpdatgeTag

type UpdatgeTag struct {
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
	tag := UpdatgeTag{}
	if err := e.Bind(&tag); err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの読み込みに失敗しました"})
	}
	randomName, err := uuid.NewRandom()
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "タグの追加に失敗しました"})
	}
	_, err = c.DB.Exec("INSERT INTO work_tags (name, description) VALUES (?, ?)", randomName.String(), "")
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "タグの追加に失敗しました"})
	}
	id := ""
	err = c.DB.QueryRow("SELECT BIN_TO_UUID(id) FROM work_tags WHERE name = ?", randomName.String()).Scan(&id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "タグの追加に失敗しました"})
	}
	err = updateTag(c.DB, id, tag)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "タグの追加に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseCreateTag{Success: true, ID: id})
}

// Update tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [put]
// @Param id path string true "tag id"
// @Param RequestUpdatgeTag body RequestUpdatgeTag true "update work"
// @Success 200 {object} ResponseUpdatgeTag
func (c Context) UpdateTag(e echo.Context) error {
	tag := UpdatgeTag{}
	if err := e.Bind(&tag); err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの読み込みに失敗しました"})
	}
	id := e.Param("id")
	err := updateTag(c.DB, id, tag)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "タグの追加に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdatgeTag{Success: true})
}

// Get tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [get]
// @Param id path string true "tag id"
// @Success 200 {object} ResponseGetTag
func (c Context) GetTag(e echo.Context) error {
	id := e.Param("id")
	tag := ResponseGetTag{}
	err := c.DB.QueryRow("SELECT name, description FROM work_tags WHERE id = UUID_TO_BIN(?)", id).Scan(&tag.Name, &tag.Description)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	return e.JSON(http.StatusOK, tag)
}

// Get tag
// @Accept json
// @Security Authorization
// @Router /work/tag/{id} [delete]
// @Param id path string true "tag id"
// @Success 200 {object} ResponseDeleteTag
func (c Context) DeleteTag(e echo.Context) error {
	id := e.Param("id")
	_, err := c.DB.Exec("DELETE FROM work_tags WHERE id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の削除に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseDeleteTag{Success: true})
}

func updateTag(db *sql.DB, id string, tag UpdatgeTag) error {
	_, err := db.Exec("UPDATE work_tags SET name = ?, description = ? WHERE id = UUID_TO_BIN(?)", tag.Name, tag.Description, id)
	if err != nil {
		return fmt.Errorf("タグの編集に失敗しました")
	}
	return nil
}
