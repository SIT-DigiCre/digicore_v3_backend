package grade_update

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetAdminGradeUpdate(ctx echo.Context, dbClient db.Client) (api.ResGetAdminGradeUpdate, *response.Error) {
	gradeUpdates := []adminGradeUpdate{}
	params := struct{}{}
	err := dbClient.Select(&gradeUpdates, "sql/grade_update/select_pending_grade_updates.sql", &params)
	if err != nil {
		return api.ResGetAdminGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "承認待ち申請一覧の取得に失敗しました", Log: err.Error()}
	}

	res := api.ResGetAdminGradeUpdate{}
	rerr := copier.Copy(&res.GradeUpdates, &gradeUpdates)
	if rerr != nil {
		return api.ResGetAdminGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "承認待ち申請一覧の取得に失敗しました", Log: rerr.Error()}
	}
	if res.GradeUpdates == nil {
		res.GradeUpdates = []api.ResGetAdminGradeUpdateObjectGradeUpdate{}
	}
	return res, nil
}
