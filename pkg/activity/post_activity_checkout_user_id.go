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

	checkOutAt := time.Now()

	executed, err := executeCheckout(dbClient, userId, requestBody.Place, checkOutAt)
	if err != nil {
		return api.BlankSuccess{}, err
	}
	if !executed {
		return api.BlankSuccess{}, &response.Error{
			Code:    http.StatusNotFound,
			Level:   "Info",
			Message: "在室中ではありません",
			Log:     "activity not found or already checked out",
		}
	}

	return api.BlankSuccess{Success: true}, nil
}
