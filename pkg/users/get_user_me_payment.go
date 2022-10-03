package users

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserMePayment(ctx echo.Context, dbClient db.Client) ([]api.ResGetUserMePayment, *response.Error) {
	userID := ctx.Get("user_id").(string)
	res, err := getUserPaymentFromUserID(userID, &dbClient)
	if err != nil {
		return []api.ResGetUserMePayment{}, err
	}
	return res, nil
}

func getUserPaymentFromUserID(userID string, dbClient db.CommonClient) ([]api.ResGetUserMePayment, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	payment := []Profile{}
	err := dbClient.Select(&payment, "sql/users/select_user_payment_from_user_id.sql", &params)
	if len(payment) == 0 {
		return []api.ResGetUserMePayment{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "支払い情報が有りません", Log: sql.ErrNoRows.Error()}
	}
	if err != nil {
		return []api.ResGetUserMePayment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return []api.ResGetUserMePayment{}, nil
}
