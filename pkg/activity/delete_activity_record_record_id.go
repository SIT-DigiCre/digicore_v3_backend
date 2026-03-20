package activity

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/admin"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteActivityRecordRecordId(ctx echo.Context, dbClient db.TransactionClient, recordId string) (api.BlankSuccess, *response.Error) {
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
				Message: "削除権限がありません",
				Log:     "ユーザーはレコードの所有者でもinfra claimを持つユーザーでもありません",
			}
		}
	}

	if err := deleteActivityRecord(dbClient, recordId); err != nil {
		return api.BlankSuccess{}, err
	}

	return api.BlankSuccess{Success: true}, nil
}

func deleteActivityRecord(dbClient db.TransactionClient, id string) *response.Error {
	params := struct {
		Id string `twowaysql:"id"`
	}{
		Id: id,
	}

	result, err := dbClient.Exec("sql/activity/delete_activity.sql", &params, false)
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "アクティビティレコードの削除に失敗しました",
			Log:     err.Error(),
		}
	}

	if err := validateDeletedRows(result); err != nil {
		return err
	}

	return nil
}

func validateDeletedRows(result sql.Result) *response.Error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "アクティビティレコードの削除結果の確認に失敗しました",
			Log:     err.Error(),
		}
	}

	if rowsAffected == 0 {
		return &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "アクティビティレコードが存在しません",
			Log:     "削除対象のアクティビティレコードが存在しないか、既に削除されています",
		}
	}

	return nil
}
