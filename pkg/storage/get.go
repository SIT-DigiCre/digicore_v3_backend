package storage

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type fileGetResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	KSize     string    `json:"k_size"`
	Extension string    `json:"extension"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}
type fileGetErrorResponse struct {
	Error string `json:"error"`
}

// Get user file url
// @Param id path string true "file id"
// @Router /storage/{fileId} [get]
// @Success 200 {object} fileGetResponse
// @Failure 404 {object} fileGetErrorResponse
// @Failure 500 {object} fileGetErrorResponse
func (c Context) GetUserFileUrl(e echo.Context) error {
	_, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fileUploadErrorResponse{Error: "ユーザ情報の読み込みに失敗しました"})
	}
	fileId := e.Param("fileId")
	data, status, err := getUserFileMetadata(c.DB, fileId)
	if err != nil {
		return e.JSON(status, fileGetErrorResponse{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, data)
}

func getUserFileMetadata(db *sql.DB, fileId string) (fileGetResponse, int, error) {
	var meta fileGetResponse
	meta.ID = fileId
	err := db.QueryRow(`SELECT BIN_TO_UUID(user_id), name, k_size, extension, is_public, created_at, updated_at FROM user_files WHERE id = UUID_TO_BIN(?)`, fileId).Scan(&meta.UserID, &meta.Name, &meta.KSize, &meta.Extension, &meta.IsPublic, &meta.CreatedAt, &meta.UpdatedAt)
	if err == sql.ErrNoRows {
		return meta, http.StatusNotFound, errors.New("指定されたファイルは見つかりませんでした")
	} else if err != nil {
		return meta, http.StatusInternalServerError, errors.New("データベースのエラーです")
	}
	url, err := getFileURL(getFileNameFromIDandExt(fileId, meta.Extension), meta.IsPublic)
	if err != nil {
		return meta, http.StatusInternalServerError, err
	}
	meta.URL = url
	return meta, http.StatusOK, nil
}
