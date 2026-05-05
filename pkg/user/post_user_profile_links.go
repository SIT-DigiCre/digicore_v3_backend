package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostUserProfileLinks(ctx echo.Context, dbClient db.TransactionClient, requestBody api.PostUserProfileLinksJSONBody) (api.ResGetUserUserId, *response.Error){
	userID := ctx.Get("user_id").(string)
	err := CreateUserProfileLinks(dbClient, userID, requestBody)
	if err != nil {
		return api.ResGetUserUserId{}, err
	}
	return GetUserUserId(ctx, dbClient, userID)
}

func CreateUserProfileLinks(dbClient db.TransactionClient, userID string, requestBody api.PostUserProfileLinksJSONBody) *response.Error {

	params := struct {
		UserID  string `twowaysql:"userId"`
		LinkUrl string `twowaysql:"linkUrl"`
	}{
		UserID:  userID,
		LinkUrl: requestBody.LinkUrl,
	}

	_, err := dbClient.Exec("sql/user/insert_user_profile_links.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "プロフィールリンクの作成に失敗しました", Log: err.Error()}
	}

	return nil
}
