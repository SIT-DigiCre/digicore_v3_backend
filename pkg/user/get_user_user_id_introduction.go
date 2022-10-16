package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserUserIDIntroduction(ctx echo.Context, dbClient db.Client, userID string) (api.ResGetUserUserIDIntroduction, *response.Error) {
	res := api.ResGetUserUserIDIntroduction{}
	introduction, err := GetUserIntroductionFromUserID(dbClient, userID)
	if err != nil {
		return api.ResGetUserUserIDIntroduction{}, err
	}
	rerr := copier.Copy(&res, &introduction)
	if rerr != nil {
		return api.ResGetUserUserIDIntroduction{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}
