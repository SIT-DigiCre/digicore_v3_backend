package admin

import (
	"encoding/json"
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

	type schoolGradeUpdate struct {
		UserId      string `json:"userId"`
		SchoolGrade int    `json:"schoolGrade"`
	}
	updates := make([]schoolGradeUpdate, 0, len(targets))

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

		updates = append(updates, schoolGradeUpdate{
			UserId:      target.UserId,
			SchoolGrade: schoolGrade + target.ApprovedGradeDiffs,
		})
	}

	updatesJson, err := json.Marshal(updates)
	if err != nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "学年の更新に失敗しました",
			Log:     err.Error(),
		}
	}

	params := struct {
		UpdatesJSON string `twowaysql:"updatesJson"`
	}{
		UpdatesJSON: string(updatesJson),
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

	return api.BlankSuccess{Success: true}, nil
}
