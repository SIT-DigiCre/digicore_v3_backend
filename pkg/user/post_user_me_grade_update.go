package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostUserMeGradeUpdate(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostUserMeGradeUpdate) (api.ResGetUserMeGradeUpdateObjectGradeUpdate, *response.Error) {
	userId := ctx.Get("user_id").(string)

	// 承認済み回数チェック
	counts := []approvedCount{}
	params := struct {
		UserId string `twowaysql:"userId"`
	}{UserId: userId}
	err := dbClient.Select(&counts, "sql/grade_update/select_approved_count_by_user_id.sql", &params)
	if err != nil {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "承認済み回数の取得に失敗しました", Log: err.Error()}
	}
	if len(counts) > 0 && counts[0].ApprovedCount >= maxApprovedCount {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "学年補正申請は最大2回までです", Log: "approved count exceeded"}
	}

	// 未処理申請の存在チェック
	pendings := []pendingCount{}
	err = dbClient.Select(&pendings, "sql/grade_update/select_pending_count_by_user_id.sql", &params)
	if err != nil {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "未処理申請の確認に失敗しました", Log: err.Error()}
	}
	if len(pendings) > 0 && pendings[0].PendingCount > 0 {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "未処理の申請が存在します。承認/却下後に再度申請してください", Log: "pending request exists"}
	}

	// 申請作成
	insertParams := struct {
		UserId    string `twowaysql:"userId"`
		GradeDiff int    `twowaysql:"gradeDiff"`
		Reason    string `twowaysql:"reason"`
	}{
		UserId:    userId,
		GradeDiff: -1,
		Reason:    requestBody.Reason,
	}
	_, rerr := dbClient.Exec("sql/grade_update/insert_grade_update.sql", &insertParams, true)
	if rerr != nil {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "申請の作成に失敗しました", Log: rerr.Error()}
	}

	// 作成した申請を取得して返却
	gradeUpdateId, rerr := dbClient.GetId()
	if rerr != nil {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "IDの取得に失敗しました", Log: rerr.Error()}
	}

	details := []gradeUpdateDetail{}
	detailParams := struct {
		GradeUpdateId string `twowaysql:"gradeUpdateId"`
	}{GradeUpdateId: gradeUpdateId}
	err = dbClient.Select(&details, "sql/grade_update/select_grade_update_by_id.sql", &detailParams)
	if err != nil {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "申請の取得に失敗しました", Log: err.Error()}
	}
	if len(details) == 0 {
		return api.ResGetUserMeGradeUpdateObjectGradeUpdate{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "申請の取得に失敗しました", Log: "created grade update not found: id=" + gradeUpdateId}
	}

	detail := details[0]
	return api.ResGetUserMeGradeUpdateObjectGradeUpdate{
		GradeUpdateId: detail.GradeUpdateId,
		GradeDiff:     detail.GradeDiff,
		Reason:        detail.Reason,
		Status:        detail.Status,
		CreatedAt:     detail.CreatedAt,
		UpdatedAt:     detail.UpdatedAt,
	}, nil
}
