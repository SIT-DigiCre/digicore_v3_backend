package work

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestCreateWork = UpdateWork

type ResponseCreateWork struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type UpdateWork struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Authers     []string `json:"authers"`
	Tags        []string `json:"tags"`
}

type RequestUpdatgeWork = UpdateWork

type ResponseUpdatgeWork struct {
	Success bool `json:"success"`
}

type Auther struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResponseGetWork struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Authers     []Auther `json:"authers"`
	Tags        []Tag    `json:"tags"`
}

type ResponseDeleteWork struct {
	Success bool `json:"success"`
}

// Create work
// @Accept json
// @Security Authorization
// @Router /work/work [post]
// @Param RequestCreateWork body RequestCreateWork true "new work"
// @Success 200 {object} ResponseCreateWork
func (c Context) CreateWork(e echo.Context) error {
	userID, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	work := UpdateWork{}
	if err := e.Bind(&work); err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの読み込みに失敗しました"})
	}
	randomName, err := uuid.NewRandom()
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の追加に失敗しました"})
	}
	_, err = c.DB.Exec("INSERT INTO works (name, description) VALUES (?, ?)", randomName.String(), "")
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の追加に失敗しました"})
	}
	id := ""
	err = c.DB.QueryRow("SELECT BIN_TO_UUID(id) FROM works WHERE name = ?", randomName.String()).Scan(&id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の追加に失敗しました"})
	}
	err = updateWork(c.DB, id, work, userID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseCreateWork{Success: true, ID: id})
}

// Update work
// @Accept json
// @Security Authorization
// @Router /work/work/{id} [put]
// @Param id path string true "work id"
// @Param RequestUpdatgeWork body RequestUpdatgeWork true "update work"
// @Success 200 {object} ResponseUpdatgeWork
func (c Context) UpdateWork(e echo.Context) error {
	userID, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	work := UpdateWork{}
	if err := e.Bind(&work); err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの読み込みに失敗しました"})
	}
	id := e.Param("id")
	err = checkWorkUpdateAuthority(c.DB, id, userID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	err = updateWork(c.DB, id, work, userID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseUpdatgeWork{Success: true})
}

// Get work
// @Accept json
// @Security Authorization
// @Router /work/work/{id} [get]
// @Param id path string true "work id"
// @Success 200 {object} ResponseGetWork
func (c Context) GetWork(e echo.Context) error {
	id := e.Param("id")
	work := ResponseGetWork{}
	err := c.DB.QueryRow("SELECT name, description FROM works WHERE id = UUID_TO_BIN(?)", id).Scan(&work.Name, &work.Description)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の取得に失敗しました"})
	}
	authers := []Auther{}
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(work_users.user_id), username FROM work_users LEFT JOIN user_profiles ON user_profiles.user_id = work_users.user_id WHERE work_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, Error{Message: "製作者の読み込みに失敗しました"})
	}
	for rows.Next() {
		auther := Auther{}
		if err := rows.Scan(&auther.ID, &auther.Name); err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "製作者の読み込みに失敗しました"})
		}
		authers = append(authers, auther)
	}
	work.Authers = authers
	tags := []Tag{}
	rows, err = c.DB.Query("SELECT BIN_TO_UUID(tag_id), name FROM work_work_tags LEFT JOIN work_tags ON work_tags.id = work_work_tags.tag_id WHERE work_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, Error{Message: "タグの読み込みに失敗しました"})
	}
	for rows.Next() {
		tag := Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "タグの読み込みに失敗しました"})
		}
		tags = append(tags, tag)
	}
	work.Tags = tags
	return e.JSON(http.StatusOK, work)
}

// Delete work
// @Accept json
// @Security Authorization
// @Router /work/work/{id} [delete]
// @Param id path string true "work id"
// @Success 200 {object} ResponseDeleteWork
func (c Context) DeleteWork(e echo.Context) error {
	userID, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	id := e.Param("id")
	err = checkWorkUpdateAuthority(c.DB, id, userID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	_, err = c.DB.Exec("DELETE FROM work_users WHERE work_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の削除に失敗しました"})
	}
	_, err = c.DB.Exec("DELETE FROM work_work_tags WHERE work_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の削除に失敗しました"})
	}
	_, err = c.DB.Exec("DELETE FROM works WHERE id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "作品の削除に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseDeleteWork{Success: true})
}

func updateWork(db *sql.DB, id string, work UpdateWork, userID string) error {
	_, err := db.Exec("UPDATE works SET name = ?, description = ? WHERE id = UUID_TO_BIN(?)", work.Name, work.Description, id)
	if err != nil {
		return fmt.Errorf("作品情報の編集に失敗しました")
	}
	work.Authers = append(work.Authers, userID)
	uniqueFlag := make(map[string]bool)
	authers := []string{}
	for _, e := range work.Authers {
		if !uniqueFlag[e] {
			uniqueFlag[e] = true
			authers = append(authers, e)
		}
	}
	_, err = db.Exec("DELETE FROM work_users WHERE work_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return fmt.Errorf("製作者の編集に失敗しました")
	}
	for _, e := range authers {
		_, err = db.Exec("INSERT INTO work_users (work_id, user_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?))", id, e)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}
	uniqueFlag = make(map[string]bool)
	tags := []string{}
	for _, e := range work.Tags {
		if !uniqueFlag[e] {
			uniqueFlag[e] = true
			tags = append(tags, e)
		}
	}
	_, err = db.Exec("DELETE FROM work_work_tags WHERE work_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	for _, e := range tags {
		_, err = db.Exec("INSERT INTO work_work_tags (work_id, tag_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?))", id, e)
		if err != nil {
			return fmt.Errorf("タグの編集に失敗しました")
		}
	}
	return nil
}

func checkWorkUpdateAuthority(db *sql.DB, workID string, userID string) error {
	tmp := ""
	err := db.QueryRow("SELECT work_id, user_id FROM work_users WHERE work_id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)", workID, userID).Scan(&tmp, &tmp)
	if err == sql.ErrNoRows {
		return fmt.Errorf("作品情報の編集の権限がありません")
	}
	if err != nil {
		return fmt.Errorf("作品情報の編集に失敗しました")
	}
	return nil
}
