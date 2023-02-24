package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/util"
	"github.com/labstack/echo/v4"
)

func PutUserMePayment(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMePayment) (api.ResGetUserMePayment, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := updateUserPayment(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetUserMePayment{}, err
	}

	return GetUserMePayment(ctx, dbClient)
}

func updateUserPayment(dbClient db.TransactionClient, userId string, requestBody api.ReqPutUserMePayment) *response.Error {
	params := struct {
		UserId       string `twowaysql:"userId"`
		Year         int    `twowaysql:"year"`
		TransferName string `twowaysql:"transferName"`
	}{
		UserId:       userId,
		Year:         util.GetSchoolYear(),
		TransferName: requestBody.TransferName,
	}
	_, err := dbClient.DuplicateUpdate("sql/user/insert_user_payment.sql", "sql/user/update_user_payment.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
