package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMeGraduated(ctx echo.Context, dbClient db.TransactionClient) (api.ResGetUserMe, *response.Error) {
	userId := ctx.Get("user_id").(string)
	profile, err := GetUserProfileFromUserId(dbClient, userId)
	if err != nil {
		return api.ResGetUserMe{}, err
	}

	if profile.SchoolGrade < 4 {
		return api.ResGetUserMe{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "4年生以上のみ卒業済みにできます",
			Log:     "school_grade is less than 4",
		}
	}

	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	_, execErr := dbClient.Exec("sql/user/update_user_is_graduated.sql", &params, false)
	if execErr != nil {
		return api.ResGetUserMe{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     execErr.Error(),
		}
	}

	return GetUserMe(ctx, dbClient)
}
