package users

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserMe(ctx echo.Context, dbClient db.Client) (api.ResGetUserMe, *response.Error) {
	userID := ctx.Get("user_id").(string)
	profile, err := GetUserProfileFromUserID(userID, dbClient)
	if err != nil {
		return api.ResGetUserMe{}, err
	}
	res := api.ResGetUserMe{}
	rerr := copier.Copy(&res, &profile)
	if rerr != nil {
		return api.ResGetUserMe{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "プロフィールの読み込みに失敗しました", Log: rerr.Error()}
	}
	return res, nil
}
