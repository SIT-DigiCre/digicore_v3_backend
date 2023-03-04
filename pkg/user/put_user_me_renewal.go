package user

import (
	"net/http"
	"strconv"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutUserMeRenewal(ctx echo.Context, dbClient db.TransactionClient) (api.ResGetUserMe, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := renewalActiveLimit(dbClient, userId)
	if err != nil {
		return api.ResGetUserMe{}, err
	}

	return GetUserMe(ctx, dbClient)
}

func renewalActiveLimit(dbClient db.TransactionClient, userId string) *response.Error {
	params := struct {
		UserId      string `twowaysql:"userId"`
		ActiveLimit string `twowaysql:"activeLimit"`
	}{
		UserId:      userId,
		ActiveLimit: strconv.Itoa(utils.GetYear()) + "-06-01",
	}
	_, err := dbClient.Exec("sql/user/update_user_active_limit.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
