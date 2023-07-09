package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMeIntroduction(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMeIntroduction) (api.ResGetUserMeIntroduction, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := updateUserIntroduction(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetUserMeIntroduction{}, err
	}

	return GetUserMeIntroduction(ctx, dbClient)
}

func updateUserIntroduction(dbClient db.TransactionClient, userId string, requestBody api.ReqPutUserMeIntroduction) *response.Error {
	params := struct {
		UserId       string `twowaysql:"userId"`
		Introduction string `twowaysql:"introduction"`
	}{
		UserId:       userId,
		Introduction: requestBody.Introduction,
	}
	_, err := dbClient.Exec("sql/user/update_user_introduction.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
