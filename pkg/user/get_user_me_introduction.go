package user

import (
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserMeIntroduction(ctx echo.Context, dbClient db.Client) (api.ResGetUserMeIntroduction, *response.Error) {
	res := api.ResGetUserMeIntroduction{}
	userID := ctx.Get("user_id").(string)
	introduction, err := GetUserIntroductionFromUserID(dbClient, userID)
	if err != nil {
		return api.ResGetUserMeIntroduction{}, err
	}
	rerr := copier.Copy(&res, &introduction)
	if rerr != nil {
		return api.ResGetUserMeIntroduction{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type introduction struct {
	Introduction string `db:"introduction"`
}

func GetUserIntroductionFromUserID(dbClient db.Client, userID string) (introduction, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	introductions := []introduction{}
	err := dbClient.Select(&introductions, "sql/user/select_user_introduction_from_user_id.sql", &params)
	if err != nil {
		return introduction{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(introductions) == 0 {
		return introduction{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "自己紹介が有りません", Log: "no rows in result"}
	}
	return introductions[0], nil
}
