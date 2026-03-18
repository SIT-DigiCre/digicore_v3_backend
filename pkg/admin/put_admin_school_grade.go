package admin

import (
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutAdminSchoolGrade(_ echo.Context, dbClient db.TransactionClient) (api.BlankSuccess, *response.Error) {
	targets := []schoolGradeUpdateTarget{}
	params := struct{}{}
	err := dbClient.Select(&targets, "sql/admin/select_school_grade_update_targets.sql", &params)
	if err != nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "学年更新対象の取得に失敗しました",
			Log:     err.Error(),
		}
	}

	for _, target := range targets {
		baseGrade, calcErr := utils.CalculateSchoolGradeFromStudentNumber(target.StudentNumber)
		if calcErr != nil {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusInternalServerError,
				Level:   "Error",
				Message: "学年の計算に失敗しました",
				Log:     fmt.Sprintf("user_id=%s student_number=%s: %v", target.UserId, target.StudentNumber, calcErr),
			}
		}

		updateParams := struct {
			UserId      string `twowaysql:"userId"`
			SchoolGrade int    `twowaysql:"schoolGrade"`
		}{
			UserId:      target.UserId,
			SchoolGrade: baseGrade + target.ApprovedGradeDiffs,
		}

		_, execErr := dbClient.Exec("sql/admin/update_user_profile_school_grade.sql", &updateParams, false)
		if execErr != nil {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusInternalServerError,
				Level:   "Error",
				Message: "学年の更新に失敗しました",
				Log:     execErr.Error(),
			}
		}
	}

	return api.BlankSuccess{Success: true}, nil
}
