package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserUserID(ctx echo.Context, dbClient db.Client, userID string) (api.ResGetUserUserID, *response.Error) {
	res := api.ResGetUserUserID{}
	profile, err := GetUserProfileFromUserID(dbClient, userID)
	if err != nil {
		return api.ResGetUserUserID{}, err
	}
	rerr := copier.Copy(&res, &profile)
	if rerr != nil {
		return api.ResGetUserUserID{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}
