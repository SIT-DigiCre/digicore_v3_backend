package admin

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutAdminGradeUpdateGradeUpdateId(ctx echo.Context, dbClient db.TransactionClient, gradeUpdateId string, requestBody api.ReqPutAdminGradeUpdateGradeUpdateId) (api.BlankSuccess, *response.Error) {
	userId := ctx.Get("user_id").(string)

	// ステータス更新（楽観ロック: status='pending'の場合のみ更新）
	updateParams := struct {
		GradeUpdateId string `twowaysql:"gradeUpdateId"`
		Status        string `twowaysql:"status"`
		ApprovedBy    string `twowaysql:"approvedBy"`
	}{
		GradeUpdateId: gradeUpdateId,
		Status:        requestBody.Status,
		ApprovedBy:    userId,
	}
	result, rerr := dbClient.Exec("sql/grade_update/update_grade_update_status.sql", &updateParams, false)
	if rerr != nil {
		return api.BlankSuccess{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ステータスの更新に失敗しました", Log: rerr.Error()}
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return api.BlankSuccess{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "指定された申請が見つからないか、既に処理済みです", Log: "grade update not found or already processed"}
	}

	// 承認の場合、user_profiles.school_gradeを更新
	if requestBody.Status == "approved" {
		details := []gradeUpdateDetail{}
		detailParams := struct {
			GradeUpdateId string `twowaysql:"gradeUpdateId"`
		}{GradeUpdateId: gradeUpdateId}
		err := dbClient.Select(&details, "sql/grade_update/select_grade_update_by_id.sql", &detailParams)
		if err != nil {
			return api.BlankSuccess{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "申請の取得に失敗しました", Log: err.Error()}
		}
		if len(details) == 0 {
			return api.BlankSuccess{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "申請の取得に失敗しました", Log: "grade update not found after status update"}
		}
		detail := details[0]
		gradeParams := struct {
			UserId    string `twowaysql:"userId"`
			GradeDiff int    `twowaysql:"gradeDiff"`
		}{
			UserId:    detail.UserId,
			GradeDiff: detail.GradeDiff,
		}
		_, rerr := dbClient.Exec("sql/grade_update/update_user_profile_grade.sql", &gradeParams, false)
		if rerr != nil {
			return api.BlankSuccess{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "学年の更新に失敗しました", Log: rerr.Error()}
		}
	}

	return api.BlankSuccess{Success: true}, nil
}
