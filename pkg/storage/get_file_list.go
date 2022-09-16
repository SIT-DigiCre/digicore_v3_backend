package storage

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type fileListResponse struct {
	FileIDs []string `json:"file_ids"`
}
type fileListErrorResponse struct {
	Error string `json:"error"`
}

// Get user file list
// @Router /storage [get]
// @Success 200 {object} fileListResponse
// @Failure 404 {object} fileListErrorResponse
// @Failure 500 {object} fileListErrorResponse
func (c Context) GetUserFileList(e echo.Context) error {
	userID, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fileUploadErrorResponse{Error: "ユーザ情報の読み込みに失敗しました"})
	}
	data, status, err := getUserFileIDs(c.DB, userID)
	if err != nil {
		return e.JSON(status, fileListErrorResponse{Error: err.Error()})
	}
	return e.JSON(status, data)
}

func getUserFileIDs(db *sql.DB, userID string) (fileListResponse, int, error) {
	var res fileListResponse
	res.FileIDs = []string{}
	rows, err := db.Query("SELECT BIN_TO_UUID(id) FROM user_files WHERE user_id = UUID_TO_BIN(?)", userID)
	if err == sql.ErrNoRows {
		return res, http.StatusOK, nil
	} else if err != nil {
		return res, http.StatusInternalServerError, errors.New("データベースのエラーです")
	}

	for rows.Next() {
		var fileID string
		if err := rows.Scan(&fileID); err != nil {
			return res, http.StatusInternalServerError, errors.New("データベースのエラーです")
		}
		res.FileIDs = append(res.FileIDs, fileID)
	}
	err = rows.Err()
	if err != nil {
		return res, http.StatusInternalServerError, errors.New("データベースのエラーです")
	}
	return res, http.StatusOK, nil
}
