package users

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetUserMePayment(ctx echo.Context, dbClient db.Client) ([]api.ResGetUserMePayment, *response.Error) {
	userID := ctx.Get("user_id").(string)
	payments, err := getUserPaymentFromUserID(userID, dbClient)
	if err != nil {
		return []api.ResGetUserMePayment{}, err
	}
	res := []api.ResGetUserMePayment{}
	rerr := copier.Copy(&res, &payments)
	if rerr != nil {
		return []api.ResGetUserMePayment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "支払い情報の読み込みに失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type payment struct {
	Checked      bool   `db:"checked"`
	TransferName string `db:"transfer_name"`
	UpdatedAt    string `db:"updated_at"`
	Year         int    `db:"year"`
}

func getUserPaymentFromUserID(userID string, dbClient db.Client) ([]payment, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	payments := []payment{}
	err := dbClient.Select(&payments, "sql/users/select_user_payment_from_user_id.sql", &params)
	if err == sql.ErrNoRows {
		return []payment{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "支払い情報が有りません", Log: sql.ErrNoRows.Error()}
	}
	if err != nil {
		return []payment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	return payments, nil
}
