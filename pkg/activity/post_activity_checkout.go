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
	latest, err := selectLatestActivity(dbClient, userId, place)
	if err != nil {
		return false, err
	}
	if latest == nil || latest.CheckedOutAt != nil {
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
	if checkOutAt.Before(latest.InitialCheckedInAt) || checkOutAt.Before(latest.CheckedInAt) {
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
		Id:           latest.ID,
		CheckedOutAt: checkOutAt,
		Note:         note,
	}

	_, execErr := dbClient.Exec("sql/activity/update_activity_checkout.sql", &params, false)
	if execErr != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Info",
			Message: "DBエラーが発生しました",
			Log:     execErr.Error(),
		}
	}

	return true, nil
}
