package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutAdminSchoolGrade(_ echo.Context, dbClient db.TransactionClient) (api.BlankSuccess, *response.Error) {
	targets := []struct {
		UserId             string `db:"user_id"`
		StudentNumber      string `db:"student_number"`
		ApprovedGradeDiffs int    `db:"approved_grade_diffs"`
	}{}
	err := dbClient.Select(&targets, "sql/admin/select_school_grade_update_targets.sql", nil)
	if err != nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "学年の更新に失敗しました",
			Log:     err.Error(),
		}
	}

	for _, target := range targets {
		schoolGrade, err := utils.CalculateSchoolGradeFromStudentNumber(target.StudentNumber)
		if err != nil {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusInternalServerError,
				Level:   "Error",
				Message: "学年の計算に失敗しました",
				Log:     err.Error(),
			}
		}

		params := struct {
			UserId      string `twowaysql:"userId"`
			SchoolGrade int    `twowaysql:"schoolGrade"`
		}{
			UserId:      target.UserId,
			SchoolGrade: schoolGrade + target.ApprovedGradeDiffs,
		}
		_, err = dbClient.Exec("sql/admin/update_user_profile_school_grade.sql", &params, false)
		if err != nil {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusInternalServerError,
				Level:   "Error",
				Message: "学年の更新に失敗しました",
				Log:     err.Error(),
			}
		}
	}

	return api.BlankSuccess{Success: true}, nil
}
