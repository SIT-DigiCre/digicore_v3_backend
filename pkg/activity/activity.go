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

type ActivityRecord struct {
	ID                  string     `db:"id"`
	UserID              string     `db:"user_id"`
	Place               string     `db:"place"`
	Note                *string    `db:"note"`
	InitialCheckedInAt  time.Time  `db:"initial_checked_in_at"`
	InitialCheckedOutAt *time.Time `db:"initial_checked_out_at"`
	CheckedInAt         time.Time  `db:"checked_in_at"`
	CheckedOutAt        *time.Time `db:"checked_out_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

func PostActivityCheckin(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostActivityCheckin) (api.BlankSuccess, *response.Error) {
	userId := ctx.Get("user_id").(string)

	checkInAt := time.Now()

	latest, err := selectLatestActivity(dbClient, userId, requestBody.Place)
	if err != nil {
		return api.BlankSuccess{}, err
	}

	// すでに在室中の場合は既存レコードをチェックアウトしてから新規レコードを作成
	if latest != nil && latest.CheckedOutAt == nil {
		if err := updateActivityCheckout(dbClient, latest.ID, checkInAt); err != nil {
			return api.BlankSuccess{}, err
		}
	}

	if err := insertActivity(dbClient, userId, requestBody.Place, checkInAt); err != nil {
		return api.BlankSuccess{}, err
	}

	return api.BlankSuccess{Success: true}, nil
}

func PostActivityCheckout(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostActivityCheckout) (api.BlankSuccess, *response.Error) {
	userId := ctx.Get("user_id").(string)
	return checkoutInternal(dbClient, userId, requestBody)
}

func PostActivityCheckoutUserId(ctx echo.Context, dbClient db.TransactionClient, userId string, requestBody api.ReqPostActivityCheckout) (api.BlankSuccess, *response.Error) {
	requestUserId := ctx.Get("user_id").(string)

	isAdmin, err := admin.CheckUserIsAdmin(dbClient, requestUserId)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	if !isAdmin {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusForbidden,
			Level:   "Info",
			Message: "管理者権限がありません",
			Log:     "user is not admin",
		}
	}

	return checkoutInternal(dbClient, userId, requestBody)
}

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
			Log:     "activity record not found",
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
				Log:     "user is not owner or admin",
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

	if err := updateActivityTimes(dbClient, recordId, checkedInAt, checkedOutAt); err != nil {
		return api.BlankSuccess{}, err
	}

	return api.BlankSuccess{Success: true}, nil
}

func checkoutInternal(dbClient db.TransactionClient, userId string, requestBody api.ReqPostActivityCheckout) (api.BlankSuccess, *response.Error) {
	checkOutAt := time.Now()

	latest, err := selectLatestActivity(dbClient, userId, requestBody.Place)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	if latest == nil || latest.CheckedOutAt != nil {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "在室中ではありません",
			Log:     "activity not found or already checked out",
		}
	}

	if err := updateActivityCheckout(dbClient, latest.ID, checkOutAt); err != nil {
		return api.BlankSuccess{}, err
	}

	return api.BlankSuccess{Success: true}, nil
}

func selectLatestActivity(dbClient db.TransactionClient, userId string, place string) (*ActivityRecord, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
		Place  string `twowaysql:"place"`
	}{
		UserId: userId,
		Place:  place,
	}

	records := []ActivityRecord{}
	err := dbClient.Select(&records, "sql/activity/select_latest_activity.sql", &params)
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

func insertActivity(dbClient db.TransactionClient, userId string, place string, checkedInAt time.Time) *response.Error {
	params := struct {
		UserId             string    `twowaysql:"userId"`
		Place              string    `twowaysql:"place"`
		InitialCheckedInAt time.Time `twowaysql:"initialCheckedInAt"`
		CheckedInAt        time.Time `twowaysql:"checkedInAt"`
	}{
		UserId:             userId,
		Place:              place,
		InitialCheckedInAt: checkedInAt,
		CheckedInAt:        checkedInAt,
	}

	_, err := dbClient.Exec("sql/activity/insert_activity.sql", &params, true)
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

func updateActivityCheckout(dbClient db.TransactionClient, id string, checkedOutAt time.Time) *response.Error {
	params := struct {
		Id           string    `twowaysql:"id"`
		CheckedOutAt time.Time `twowaysql:"checkedOutAt"`
	}{
		Id:           id,
		CheckedOutAt: checkedOutAt,
	}

	_, err := dbClient.Exec("sql/activity/update_activity_checkout.sql", &params, false)
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
