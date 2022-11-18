package storage

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetStorageMyfile(ctx echo.Context, dbClient db.Client) (api.ResGetStorageMyfile, *response.Error) {
	res := api.ResGetStorageMyfile{}
	userId := ctx.Get("user_id").(string)
	files, err := getFileListFromUserID(dbClient, userId)
	if err != nil {
		return api.ResGetStorageMyfile{}, err
	}
	rerr := copier.Copy(&res.Files, &files)
	if rerr != nil {
		return api.ResGetStorageMyfile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type file struct {
	FileId    string `db:"id"`
	Name      string `db:"name"`
	Extension string `db:"extension"`
	IsPublic  bool   `db:"is_public"`
	KSize     string `db:"k_size"`
	UserId    string `db:"user_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func getFileListFromUserID(dbClient db.Client, userId string) ([]file, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	files := []file{}
	err := dbClient.Select(&files, "sql/storage/select_storage_from_user_id.sql", &params)
	if err != nil {
		return []file{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(files) == 0 {
		return []file{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "ファイルが有りません", Log: "no rows in result"}
	}
	return files, nil
}
