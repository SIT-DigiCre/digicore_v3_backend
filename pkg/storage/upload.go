package storage

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type fileUploadRequest struct {
	Name string `json:"name"`
	File string `json:"file"`
}
type fileUploadResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
type fileUploadErrorResponse struct {
	Error string `json:"error"`
}

func (f fileUploadRequest) validate() error {
	if len(f.File) == 0 {
		return errors.New("ファイルが空です")
	}
	if len(f.Name) == 0 || 255 < utf8.RuneCountInString(f.Name) {
		return errors.New("ファイル名は1文字以上255文字以下である必要があります")
	}
	return nil
}

func uploadBytes(data []byte, fileId string, extension string) (int, error) {
	if env.ConohaFileUploadMaxSize < len(data) {
		return http.StatusBadRequest, errors.New("ファイルサイズが大きすぎます")
	}
	token, err := GetToken()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	url := getFileURL(fileId, extension)
	_, err = httpRequest(http.MethodPut, url, bytes.NewReader(data), &token)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func createUserFile(db *sql.DB, userId string, data []byte, fileName string) (string, int, error) {
	id := uuid.New().String()
	extension := getExtension(fileName)
	md5Hash := fmt.Sprintf("%x", md5.Sum(data))
	duplicateFileName := ""
	err := db.QueryRow(`SELECT name FROM user_files WHERE user_id = UUID_TO_BIN(?) AND md5_hash = ?`, userId, md5Hash).Scan(&duplicateFileName)
	// エラーがない = 該当するレコードが存在した場合
	if err == nil {
		return "", http.StatusBadRequest, errors.New(fmt.Sprintf("アップロードされたファイルは既に%sという名前でアップロードされています", duplicateFileName))
	}
	if err != sql.ErrNoRows {
		return "", http.StatusInternalServerError, errors.New("データベースのエラーです")
	}
	status, err := uploadBytes(data, id, extension)
	if err != nil {
		return "", status, err
	}
	kSize := len(data) / 1024
	_, err = db.Exec(
		`INSERT INTO user_files (id, user_id, name, k_size, md5_hash, extension) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?)`,
		id,
		userId,
		fileName,
		kSize,
		md5Hash,
		extension,
	)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("データベースのエラーです")
	}
	return id, http.StatusCreated, nil
}

// Upload user data
// @Accept json
// @Param fileUploadRequest body fileUploadRequest true "base64 encoded file and file name"
// @Router /storage [post]
// @Success 200 {object} fileUploadResponse
func (c Context) UploadUserfile(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fileUploadErrorResponse{Error: "ユーザ情報の読み込みに失敗しました"})
	}
	var fileUpload fileUploadRequest
	if err := e.Bind(&fileUpload); err != nil {
		return e.JSON(http.StatusBadRequest, fileUploadErrorResponse{Error: "データの読み込みに失敗しました"})
	}
	if err := fileUpload.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, fileUploadErrorResponse{Error: err.Error()})
	}
	data, err := base64.StdEncoding.DecodeString(fileUpload.File)
	if err != nil {
		return e.JSON(http.StatusBadRequest, fileUploadErrorResponse{Error: "ファイルのデコードに失敗しました"})
	}
	fileId, status, err := createUserFile(c.DB, userId, data, fileUpload.Name)
	if err != nil {
		return e.JSON(status, fileUploadErrorResponse{Error: err.Error()})
	}

	url := getFileURL(fileId, getExtension(fileUpload.Name))
	return e.JSON(http.StatusCreated, fileUploadResponse{ID: fileId, Name: fileUpload.Name, URL: url})
}
