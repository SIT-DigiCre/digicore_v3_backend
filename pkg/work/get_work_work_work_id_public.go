package work

import (
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/storage"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetWorkWorkWorkIdPublic(ctx echo.Context, dbClient db.Client, workId string) (api.ResGetWorkWorkWorkIdPublic, *response.Error) {
	res := api.ResGetWorkWorkWorkIdPublic{}
	tmp_res, err := GetWorkWorkWorkId(ctx, dbClient, workId)
	if err != nil {
		return api.ResGetWorkWorkWorkIdPublic{}, err
	}
	rerr := copier.Copy(&res, &tmp_res)
	if rerr != nil {
		return api.ResGetWorkWorkWorkIdPublic{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	viewExtention := []string{"jpg", "jpeg", "png", "gif"}
	if 0 < len(tmp_res.Files) {
		for _, file := range tmp_res.Files {
			viewAble := false
			for _, extention := range viewExtention {
				if strings.HasSuffix(file.Name, extention) {
					viewAble = true
				}
			}
			if viewAble {
				file, _ := storage.GetStorageFileId(ctx, dbClient, file.FileId)
				res.FileUrl = &file.Url
				res.FileName = &file.Name
				break
			}
		}
	}
	return res, nil
}
