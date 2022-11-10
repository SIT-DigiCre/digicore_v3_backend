package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserUserId(ctx echo.Context, dbClient db.Client, userId string) (api.ResGetUserUserId, *response.Error) {
	res := api.ResGetUserUserId{}
	profile, err := GetUserProfileFromUserId(dbClient, userId)
	if err != nil {
		return api.ResGetUserUserId{}, err
	}
	rerr := copier.Copy(&res, &profile)
	if rerr != nil {
		return api.ResGetUserUserId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}
