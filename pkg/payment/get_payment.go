package payment

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetPayment(ctx echo.Context, dbClient db.Client, params api.GetPaymentParams) (api.ResGetPayment, *response.Error) {
	res := api.ResGetPayment{}
	payments, err := getPaymentList(dbClient, params.Year)
	if err != nil {
		return api.ResGetPayment{}, err
	}
	rerr := copier.Copy(&res.Payments, &payments)
	if rerr != nil {
		return api.ResGetPayment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	if res.Payments == nil {
		res.Payments = []api.ResGetPaymentObjectPayment{}
	}
	return res, nil
}

func getPaymentList(dbClient db.Client, year *int) ([]payment, *response.Error) {
	searchYear := utils.GetSchoolYear()
	if year != nil {
		searchYear = *year
	}
	params := struct {
		Year int `twowaysql:"year"`
	}{
		Year: searchYear,
	}
	payments := []payment{}
	err := dbClient.Select(&payments, "sql/payment/select_payment.sql", &params)
	if err != nil {
		return []payment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return payments, nil
}
