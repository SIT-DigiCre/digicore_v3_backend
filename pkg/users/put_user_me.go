package users

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMe(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMe) (api.ResGetUserMe, *response.Error) {
	userID := ctx.Get("user_id").(string)
	err := updateUserProfile(dbClient, userID, requestBody)
	if err != nil {
		return api.ResGetUserMe{}, err
	}

	profile, err := getUserProfileFromUserID(userID, dbClient)
	if err != nil {
		return api.ResGetUserMe{}, err
	}
	res := api.ResGetUserMe{
		ActiveLimit:           profile.ActiveLimit,
		DiscordUserid:         profile.DiscordUserID,
		IconUrl:               profile.IconUrl,
		SchoolGrade:           profile.SchoolGrade,
		ShortSelfIntroduction: profile.ShortSelfIntroduction,
		StudentNumber:         profile.StudentNumber,
		UserId:                profile.UserId,
		Username:              profile.Username,
	}
	return res, nil
}

func updateUserProfile(dbClient db.TransactionClient, userID string, requestBody api.ReqPutUserMe) *response.Error {
	params := struct {
		UserID                string `twowaysql:"userID"`
		IconUrl               string `twowaysql:"iconURL"`
		SchoolGrade           int    `twowaysql:"schoolGrade"`
		ShortSelfIntroduction string `twowaysql:"shortSelfIntroduction"`
		Username              string `twowaysql:"username"`
	}{
		UserID:                userID,
		IconUrl:               requestBody.IconUrl,
		SchoolGrade:           requestBody.SchoolGrade,
		ShortSelfIntroduction: requestBody.ShortSelfIntroduction,
		Username:              requestBody.Username,
	}
	_, err := dbClient.DuplicateUpdate("sql/users/insert_user_profile.sql", "sql/users/update_user_profile.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return nil
}
