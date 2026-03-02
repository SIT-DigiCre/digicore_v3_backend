package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserMeGradeUpdate(ctx echo.Context, dbClient db.Client) (api.ResGetUserMeGradeUpdate, *response.Error) {
	userId := ctx.Get("user_id").(string)

	params := struct {
		UserId string `twowaysql:"userId"`
	}{UserId: userId}
	gradeUpdates := []gradeUpdate{}
	err := dbClient.Select(&gradeUpdates, "sql/grade_update/select_grade_updates_by_user_id.sql", &params)
	if err != nil {
		return api.ResGetUserMeGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "学年補正申請一覧の取得に失敗しました", Log: err.Error()}
	}

	res := api.ResGetUserMeGradeUpdate{}
	rerr := copier.Copy(&res.GradeUpdates, &gradeUpdates)
	if rerr != nil {
		return api.ResGetUserMeGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "学年補正申請一覧の取得に失敗しました", Log: rerr.Error()}
	}
	if res.GradeUpdates == nil {
		res.GradeUpdates = []api.ResGetUserMeGradeUpdateObjectGradeUpdate{}
	}
	return res, nil
}
