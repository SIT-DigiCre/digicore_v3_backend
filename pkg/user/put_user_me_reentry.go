package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutUserMeReentry(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMeReentry) (api.ResGetUserMeReentryObjectReentry, *response.Error) {
	userId := ctx.Get("user_id").(string)

	profile, profileErr := GetUserProfileFromUserId(dbClient, userId)
	if profileErr != nil {
		return api.ResGetUserMeReentryObjectReentry{}, profileErr
	}
	if profile.IsMember {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "有効なアカウントでは再入部申請できません", Log: "member user requested reentry"}
	}

	params := struct {
		UserId string `twowaysql:"userId"`
	}{UserId: userId}

	lockRows := []struct {
		UserId string `db:"user_id"`
	}{}
	err := dbClient.Select(&lockRows, "sql/reentry/select_user_profile_for_update_by_user_id.sql", &params)
	if err != nil {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請の確認に失敗しました", Log: err.Error()}
	}
	if len(lockRows) == 0 {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "プロフィールが有りません", Log: "user profile not found while locking reentry"}
	}

	pendings := []reentryPendingCount{}
	err = dbClient.Select(&pendings, "sql/reentry/select_pending_reentry_count_by_user_id.sql", &params)
	if err != nil {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "未処理申請の確認に失敗しました", Log: err.Error()}
	}
	if len(pendings) > 0 && pendings[0].PendingCount > 0 {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "再入部申請の確認中です。振込案内ページをご確認ください", Log: "pending reentry exists"}
	}

	counts := []reentryTotalCount{}
	err = dbClient.Select(&counts, "sql/reentry/select_reentry_total_count_by_user_id.sql", &params)
	if err != nil {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請回数の取得に失敗しました", Log: err.Error()}
	}
	totalCount := 0
	if len(counts) > 0 {
		totalCount = counts[0].TotalCount
	}
	if totalCount >= maxReentryCount {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "再入部申請は最大2回までです", Log: "reentry count exceeded"}
	}

	_, errRes := updateUserPayment(dbClient, userId, api.ReqPutUserMePayment(requestBody))
	if errRes != nil {
		return api.ResGetUserMeReentryObjectReentry{}, errRes
	}

	insertParams := struct {
		UserId string `twowaysql:"userId"`
	}{UserId: userId}
	_, execErr := dbClient.Exec("sql/reentry/insert_reentry.sql", &insertParams, true)
	if execErr != nil {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請の作成に失敗しました", Log: execErr.Error()}
	}

	reentryId, idErr := dbClient.GetId()
	if idErr != nil {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請IDの取得に失敗しました", Log: idErr.Error()}
	}

	details := []reentryDetail{}
	detailParams := struct {
		ReentryId string `twowaysql:"reentryId"`
	}{ReentryId: reentryId}
	err = dbClient.Select(&details, "sql/reentry/select_reentry_by_id.sql", &detailParams)
	if err != nil {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請の取得に失敗しました", Log: err.Error()}
	}
	if len(details) == 0 {
		return api.ResGetUserMeReentryObjectReentry{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "再入部申請の取得に失敗しました", Log: "created reentry not found: id=" + reentryId}
	}

	notifyReentryApplied(profile.StudentNumber)

	detail := details[0]
	return api.ResGetUserMeReentryObjectReentry{
		ReentryId: detail.ReentryId,
		Status:    detail.Status,
		CreatedAt: detail.CreatedAt,
		UpdatedAt: detail.UpdatedAt,
	}, nil
}
