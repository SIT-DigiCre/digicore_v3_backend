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
func PutActivityRecordRecordId(ctx echo.Context, dbClient db.TransactionClient, recordId string, requestBody ActivityRecordUpdateRequest) (api.BlankSuccess, *response.Error) {
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
		hasInfra, aerr := admin.CheckUserHasClaim(dbClient, requestUserId, "infra")
		if aerr != nil {
			return api.BlankSuccess{}, aerr
		}
		if !hasInfra {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusForbidden,
				Level:   "Info",
				Message: "編集権限がありません",
				Log:     "ユーザーはレコードの所有者でもinfra claimを持つユーザーでもありません",
			}
		}
	}

	checkedInAt := record.CheckedInAt
	checkedOutAt := record.CheckedOutAt

	switch requestBody.ActivityType {
	case "checkin":
		if requestBody.Time.IsZero() {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "更新時刻が指定されていません",
				Log:     "checkin更新でtimeがnilです",
			}
		}
		checkedInAt = requestBody.Time
		if checkedOutAt != nil && checkedOutAt.Before(checkedInAt) {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "チェックアウト時刻はチェックイン時刻より後である必要があります",
				Log:     "チェックアウト時刻がチェックイン時刻より前です",
			}
		}
	case "checkout":
		if requestBody.Time.IsZero() {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "更新時刻が指定されていません",
				Log:     "checkout更新でtimeがnilです",
			}
		}
		if requestBody.Time.Before(record.InitialCheckedInAt) || requestBody.Time.Before(record.CheckedInAt) {
			return api.BlankSuccess{}, &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "チェックアウト時刻はチェックイン時刻より後である必要があります",
				Log:     "チェックアウト時刻が初回または最新のチェックイン時刻より前です",
			}
		}
		checkedOutAt = &requestBody.Time
	default:
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "activity_typeが不正です",
			Log:     "許可されていないactivity_typeです",
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
