package work

import (
	"encoding/base64"
	"net/http"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const maxThumbnailBytes = 5 * 1024 * 1024

func thumbnailError(code int, message string) *response.Error {
	return &response.Error{Code: code, Level: "Info", Message: message, Log: message}
}

func decodeThumbnail(name, encoded string) ([]byte, string, *response.Error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, "", thumbnailError(http.StatusBadRequest, "画像のBase64形式が不正です")
	}
	if len(data) > maxThumbnailBytes {
		return nil, "", thumbnailError(http.StatusRequestEntityTooLarge, "サムネイルは5MiB以下にしてください")
	}
	header := data
	if len(header) > 512 {
		header = header[:512]
	}
	contentType := http.DetectContentType(header)
	expected := map[string]string{".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png", ".gif": "image/gif"}[strings.ToLower(filepath.Ext(name))]
	if expected == "" || contentType != expected {
		return nil, "", thumbnailError(http.StatusBadRequest, "JPEG・PNG・GIF画像を、内容と一致する拡張子で指定してください")
	}
	return data, contentType, nil
}

func PostWorkThumbnail(ctx echo.Context, client db.TransactionClient, body api.ReqPostWorkThumbnail) (workObjectThumbnail, *response.Error) {
	userID := ctx.Get("user_id").(string)
	if len(body.Name) == 0 || utf8.RuneCountInString(body.Name) > 255 {
		return workObjectThumbnail{}, thumbnailError(http.StatusBadRequest, "ファイル名は1〜255文字で指定してください")
	}
	data, contentType, err := decodeThumbnail(body.Name, body.File)
	if err != nil {
		return workObjectThumbnail{}, err
	}
	id := uuid.NewString()
	key := "work-thumbnails/" + id + strings.ToLower(filepath.Ext(body.Name))
	rerr := storage.UploadWorkThumbnail(data, key, contentType)
	if rerr != nil {
		return workObjectThumbnail{}, &response.Error{Code: 500, Level: "Error", Message: "サムネイルの保存に失敗しました", Log: rerr.Error()}
	}
	params := struct {
		Id          string `twowaysql:"id"`
		UserId      string `twowaysql:"userId"`
		StorageKey  string `twowaysql:"storageKey"`
		ContentType string `twowaysql:"contentType"`
	}{id, userID, key, contentType}
	_, rerr = client.Exec("sql/work/insert_work_thumbnail.sql", &params, false)
	if rerr != nil {
		return workObjectThumbnail{}, &response.Error{Code: 500, Level: "Error", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return workObjectThumbnail{ThumbnailId: id, Url: storage.WorkThumbnailURL(key)}, nil
}
