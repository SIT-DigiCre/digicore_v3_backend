package users

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserMePayment(ctx echo.Context, dbClient db.Client) (api.ResGetUserMePayment, *response.Error) {
	res := api.ResGetUserMePayment{}
	userID := ctx.Get("user_id").(string)
	history, err := getUserPaymentFromUserID(dbClient, userID)
	if err != nil {
		return api.ResGetUserMePayment{}, err
	}
	rerr := copier.Copy(&res.History, &history)
	if rerr != nil {
		return api.ResGetUserMePayment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type payment struct {
	Checked      bool   `db:"checked"`
	TransferName string `db:"transfer_name"`
	UpdatedAt    string `db:"updated_at"`
	Year         int    `db:"year"`
}

func getUserPaymentFromUserID(dbClient db.Client, userID string) ([]payment, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	payments := []payment{}
	err := dbClient.Select(&payments, "sql/users/select_user_payment_from_user_id.sql", &params)
	if err != nil {
		return []payment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(payments) == 0 {
		return []payment{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "支払い情報が有りません", Log: "no rows in result"}
	}
	return payments, nil
}
