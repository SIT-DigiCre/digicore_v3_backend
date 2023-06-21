package user

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutUserMePayment(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMePayment) (api.ResGetUserMePayment, *response.Error) {
	userId := ctx.Get("user_id").(string)
	update, err := updateUserPayment(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetUserMePayment{}, err
	}
	profile, _ := GetUserProfileFromUserId(dbClient, userId)
	utils.NoticeMattermost(fmt.Sprintf("%s(%s)が振込申請を行いました", profile.StudentNumber, userId), "digicore-notice", "digicore-notice", "bell")
	if update {
		return GetUserMePayment(ctx, dbClient)
	}
	err = utils.RenewalActiveLimit(dbClient, userId, utils.GetAfterDate(0, 1, 0))
	if err != nil {
		return api.ResGetUserMePayment{}, err
	}
	return GetUserMePayment(ctx, dbClient)
}

func updateUserPayment(dbClient db.TransactionClient, userId string, requestBody api.ReqPutUserMePayment) (bool, *response.Error) {
	params := struct {
		UserId       string `twowaysql:"userId"`
		Year         int    `twowaysql:"year"`
		TransferName string `twowaysql:"transferName"`
	}{
		UserId:       userId,
		Year:         utils.GetSchoolYear(),
		TransferName: requestBody.TransferName,
	}
	_, update, err := dbClient.DuplicateUpdate("sql/user/insert_user_payment.sql", "sql/user/update_user_payment.sql", &params)
	if err != nil {
		return false, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return update, nil
}
