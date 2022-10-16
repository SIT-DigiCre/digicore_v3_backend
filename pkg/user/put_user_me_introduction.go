package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMeIntroduction(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMeIntroduction) (api.ResGetUserMeIntroduction, *response.Error) {
	userID := ctx.Get("user_id").(string)
	err := updateUserIntroduction(dbClient, userID, requestBody)
	if err != nil {
		return api.ResGetUserMeIntroduction{}, err
	}

	return GetUserMeIntroduction(ctx, dbClient)
}

func updateUserIntroduction(dbClient db.TransactionClient, userID string, requestBody api.ReqPutUserMeIntroduction) *response.Error {
	params := struct {
		UserID       string `twowaysql:"userID"`
		Introduction string `twowaysql:"introduction"`
	}{
		UserID:       userID,
		Introduction: requestBody.Introduction,
	}
	_, err := dbClient.Exec("sql/user/update_user_introduction.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
