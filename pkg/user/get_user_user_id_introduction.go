package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserUserIdIntroduction(ctx echo.Context, dbClient db.Client, userId string) (api.ResGetUserUserIdIntroduction, *response.Error) {
	res := api.ResGetUserUserIdIntroduction{}
	introduction, err := GetUserIntroductionFromUserId(dbClient, userId)
	if err != nil {
		return api.ResGetUserUserIdIntroduction{}, err
	}
	rerr := copier.Copy(&res, &introduction)
	if rerr != nil {
		return api.ResGetUserUserIdIntroduction{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}
