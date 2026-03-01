package activity

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PostActivityCheckoutUserId(ctx echo.Context, dbClient db.TransactionClient, userId string, requestBody api.ReqPostActivityCheckout) (api.BlankSuccess, *response.Error) {
	note := "管理者による退室"
	executed, err := executeCheckout(dbClient, userId, requestBody.Place, requestBody.CheckoutAt, &note)
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
