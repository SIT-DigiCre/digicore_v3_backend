package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMePrivate(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMePrivate) (api.ResGetUserMePrivate, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := updateUserPrivate(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetUserMePrivate{}, err
	}

	return GetUserMePrivate(ctx, dbClient)
}

func updateUserPrivate(dbClient db.TransactionClient, userId string, requestBody api.ReqPutUserMePrivate) *response.Error {
	params := struct {
		UserId                string  `twowaysql:"userId"`
		FirstName             string  `twowaysql:"firstName"`
		LastName              string  `twowaysql:"lastName"`
		FirstNameKana         string  `twowaysql:"firstNameKana"`
		LastNameKana          string  `twowaysql:"lastNameKana"`
		IsMale                bool    `twowaysql:"isMale"`
		PhoneNumber           string  `twowaysql:"phoneNumber"`
		Address               string  `twowaysql:"address"`
		ParentName            string  `twowaysql:"parentName"`
		ParentCellphoneNumber string  `twowaysql:"parentCellphoneNumber"`
		ParentHomephoneNumber *string `twowaysql:"parentHomephoneNumber"`
		ParentAddress         string  `twowaysql:"parentAddress"`
	}{
		UserId:                userId,
		FirstName:             requestBody.FirstName,
		LastName:              requestBody.LastName,
		FirstNameKana:         requestBody.FirstNameKana,
		LastNameKana:          requestBody.LastNameKana,
		IsMale:                requestBody.IsMale,
		PhoneNumber:           requestBody.PhoneNumber,
		Address:               requestBody.Address,
		ParentName:            requestBody.ParentName,
		ParentCellphoneNumber: requestBody.ParentCellphoneNumber,
		ParentHomephoneNumber: requestBody.ParentHomephoneNumber,
		ParentAddress:         requestBody.ParentAddress,
	}
	_, err := dbClient.DuplicateUpdate("sql/user/insert_user_private.sql", "sql/user/update_user_private.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
