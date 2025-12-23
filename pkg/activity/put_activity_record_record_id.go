package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

// 在室レコードの入退室時刻更新
func PutActivityRecordRecordId(ctx echo.Context, dbClient db.TransactionClient, recordId string, requestBody api.ReqPutActivityRecordRecordId) (api.BlankSuccess, *response.Error) {
	record, err := selectActivityFromId(dbClient, recordId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	if record == nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "アクティビティレコードが存在しません",
			Log:     "アクティビティレコードが見つかりません",
		}
	}

	requestUserId := ctx.Get("user_id").(string)
	if requestUserId != record.UserID {
		isAdmin, aerr := admin.CheckUserIsAdmin(dbClient, requestUserId)
		if aerr != nil {
			return api.BlankSuccess{}, aerr
		}
		if !isAdmin {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusForbidden,
				Level:   "Info",
				Message: "編集権限がありません",
				Log:     "ユーザーはレコードの所有者でも管理者でもありません",
			}
		}
	}

	checkedInAt := record.CheckedInAt
	checkedOutAt := record.CheckedOutAt

	if requestBody.CheckedInAt != nil {
		checkedInAt = *requestBody.CheckedInAt
	}
	if requestBody.CheckedOutAt != nil {
		checkedOutAt = requestBody.CheckedOutAt
	}

	// チェックアウト時刻がチェックイン時刻より前でないことを確認
	if checkedOutAt != nil && checkedOutAt.Before(checkedInAt) {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "チェックアウト時刻はチェックイン時刻より後である必要があります",
			Log:     "チェックアウト時刻がチェックイン時刻より前です",
		}
	}

	if err := updateActivityTimes(dbClient, recordId, checkedInAt, checkedOutAt); err != nil {
		return api.BlankSuccess{}, err
	}

	return api.BlankSuccess{Success: true}, nil
}

func selectActivityFromId(dbClient db.TransactionClient, id string) (*ActivityRecord, *response.Error) {
	params := struct {
		Id string `twowaysql:"id"`
	}{
		Id: id,
	}

	records := []ActivityRecord{}
	err := dbClient.Select(&records, "sql/activity/select_activity_from_id.sql", &params)
	if err != nil {
		return nil, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}
	if len(records) == 0 {
		return nil, nil
	}
	return &records[0], nil
}

func updateActivityTimes(dbClient db.TransactionClient, id string, checkedInAt time.Time, checkedOutAt *time.Time) *response.Error {
	params := struct {
		Id           string     `twowaysql:"id"`
		CheckedInAt  time.Time  `twowaysql:"checkedInAt"`
		CheckedOutAt *time.Time `twowaysql:"checkedOutAt"`
	}{
		Id:           id,
		CheckedInAt:  checkedInAt,
		CheckedOutAt: checkedOutAt,
	}

	_, err := dbClient.Exec("sql/activity/update_activity_times.sql", &params, false)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     err.Error(),
		}
	}
	return nil
}
