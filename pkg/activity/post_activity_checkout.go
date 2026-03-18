package activity

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostActivityCheckout(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPostActivityCheckout) (api.BlankSuccess, *response.Error) {
	userId := ctx.Get("user_id").(string)

	if err := lockActivityUser(dbClient, userId); err != nil {
		return api.BlankSuccess{}, err
	}

	executed, err := executeCheckout(dbClient, userId, requestBody.Place, requestBody.CheckoutAt, nil)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	if !executed {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "在室中ではありません",
			Log:     "アクティビティレコードが見つからないか、既にチェックアウト済みです",
		}
	}

	return api.BlankSuccess{Success: true}, nil
}

func executeCheckout(dbClient db.TransactionClient, userId string, place string, requestedCheckoutAt *time.Time, note *string) (executed bool, err *response.Error) {
	current, err := selectCurrentActivity(dbClient, userId, place)
	if err != nil {
		return false, err
	}
	if current == nil {
		return false, nil
	}

	checkOutAt := time.Now()
	if requestedCheckoutAt != nil {
		if requestedCheckoutAt.IsZero() {
			return false, &response.Error{
				Code:    http.StatusBadRequest,
				Level:   "Info",
				Message: "チェックアウト時刻が不正です",
				Log:     "checkout_atがゼロ値です",
			}
		}
		checkOutAt = *requestedCheckoutAt
	}
	if checkOutAt.Before(current.InitialCheckedInAt) || checkOutAt.Before(current.CheckedInAt) {
		return false, &response.Error{
			Code:    http.StatusBadRequest,
			Level:   "Info",
			Message: "チェックアウト時刻はチェックイン時刻以降である必要があります",
			Log:     "指定されたチェックアウト時刻がチェックイン時刻より前です",
		}
	}

	params := struct {
		Id           string    `twowaysql:"id"`
		CheckedOutAt time.Time `twowaysql:"checkedOutAt"`
		Note         *string   `twowaysql:"note"`
	}{
		Id:           current.ID,
		CheckedOutAt: checkOutAt,
		Note:         note,
	}

	result, execErr := dbClient.Exec("sql/activity/update_activity_checkout.sql", &params, false)
	if execErr != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     execErr.Error(),
		}
	}

	rowsAffected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     rowsErr.Error(),
		}
	}
	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
