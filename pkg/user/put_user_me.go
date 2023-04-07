package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMe(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMe) (api.ResGetUserMe, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := UpdateUserProfile(dbClient, userId, requestBody)
	if err != nil {
		return api.ResGetUserMe{}, err
	}

	return GetUserMe(ctx, dbClient)
}

func UpdateUserProfile(dbClient db.TransactionClient, userId string, requestBody api.ReqPutUserMe) *response.Error {
	params := struct {
		UserId            string `twowaysql:"userId"`
		IconUrl           string `twowaysql:"iconUrl"`
		SchoolGrade       int    `twowaysql:"schoolGrade"`
		ShortIntroduction string `twowaysql:"shortIntroduction"`
		Username          string `twowaysql:"username"`
	}{
		UserId:            userId,
		IconUrl:           requestBody.IconUrl,
		SchoolGrade:       requestBody.SchoolGrade,
		ShortIntroduction: requestBody.ShortIntroduction,
		Username:          requestBody.Username,
	}
	_, _, err := dbClient.DuplicateUpdate("sql/user/insert_user_profile.sql", "sql/user/update_user_profile.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
