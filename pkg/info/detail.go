package info

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseGetInfo struct {
	ID        string    `json:"id"`
	Auther    Auther    `json:"auther"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Auther struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

type RequestPostInfo struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// Get info list
// @Accept json
// @Security Authorization
// @Router /info/{id} [get]
// @Param id path string true "info id"
// @Success 200 {object} ResponseGetInfo
func (c Context) GetInfo(e echo.Context) error {
	id := e.Param("id")
	res, err := GetInfoFromDB(c.DB, id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}
	return e.JSON(http.StatusOK, res)
}

// Post info
// @Accept json
// @Security Authorization
// @Router /info [post]
// @Param RequestPostInfo body RequestPostInfo true "update info"
// @Success 200 {object} ResponseGetInfo
func (c Context) PostInfo(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "ユーザーIDの取得に失敗しました"})
	}

	info := ResponseGetInfo{}
	if err := e.Bind(&info); err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの追加に失敗しました"})
	}
	uuidObj, err := uuid.NewRandom()
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの追加に失敗しました"})
	}
	id := uuidObj.String()
	_, err = c.DB.Exec("INSERT INTO info_boards (id, title, body, user_id) VALUE (UUID_TO_BIN(?), ?, ?, UUID_TO_BIN(?))", id, info.Title, info.Body, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの追加に失敗しました"})
	}

	res, err := GetInfoFromDB(c.DB, id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

// Update info
// @Accept json
// @Security Authorization
// @Router /info/{id} [put]
// @Param id path string true "info id"
// @Param RequestPostInfo body RequestPostInfo true "update info"
// @Success 200 {object} ResponseGetInfo
func (c Context) UpdateInfo(e echo.Context) error {
	id := e.Param("id")
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "ユーザーIDの取得に失敗しました"})
	}

	info := ResponseGetInfo{}
	if err := e.Bind(&info); err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの更新に失敗しました"})
	}
	_, err = c.DB.Exec("UPDATE info_boards SET title = ?, body = ? WHERE id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)", info.Title, info.Body, id, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: "データの更新に失敗しました"})
	}

	res, err := GetInfoFromDB(c.DB, id)
	if err != nil {
		return e.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	return e.JSON(http.StatusOK, res)
}

func GetInfoFromDB(db *sql.DB, id string) (ResponseGetInfo, error) {
	res := ResponseGetInfo{}
	err := db.QueryRow("SELECT BIN_TO_UUID(info_boards.id), title, body, updated_at, created_at, BIN_TO_UUID(user_profiles.user_id), user_profiles.username, user_profiles.icon_url FROM info_boards LEFT JOIN user_profiles ON user_profiles.user_id = info_boards.user_id WHERE info_boards.id = UUID_TO_BIN(?)", id).
		Scan(&res.ID, &res.Title, &res.Body, &res.UpdatedAt, &res.CreatedAt, &res.Auther.ID, &res.Auther.Name, &res.Auther.IconURL)
	if err == sql.ErrNoRows {
		return ResponseGetInfo{}, fmt.Errorf("データがありません")
	} else if err != nil {
		return ResponseGetInfo{}, fmt.Errorf("データの取得に失敗しました")
	}
	return res, nil
}
