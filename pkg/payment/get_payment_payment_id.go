package payment

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetPaymentPaymentId(ctx echo.Context, dbClient db.Client, paymentId string) (api.ResGetPaymentPaymentId, *response.Error) {
	res := api.ResGetPaymentPaymentId{}
	payment, err := getPaymentFromPaymentId(dbClient, paymentId)
	if err != nil {
		return api.ResGetPaymentPaymentId{}, err
	}
	rerr := copier.Copy(&res, &payment)
	if rerr != nil {
		return api.ResGetPaymentPaymentId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

type payment struct {
	Checked       bool   `db:"checked"`
	PaymentId     string `db:"id"`
	StudentNumber string `db:"student_number"`
	TransferName  string `db:"transfer_name"`
	UserId        string `db:"user_id"`
	Note          string `db:"note"`
}

func getPaymentFromPaymentId(dbClient db.Client, paymentId string) (payment, *response.Error) {
	params := struct {
		PaymentId string `twowaysql:"paymentId"`
	}{
		PaymentId: paymentId,
	}
	payments := []payment{}
	err := dbClient.Select(&payments, "sql/payment/select_payment_from_payment_id.sql", &params)
	if err != nil {
		return payment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(payments) == 0 {
		return payment{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "支払い情報が有りません", Log: "no rows in result"}
	}
	return payments[0], nil
}
