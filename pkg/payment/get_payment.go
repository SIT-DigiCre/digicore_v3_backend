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
	payments, err := getPaymentList(dbClient, params.Offset, params.Year)
	if err != nil {
		return api.ResGetPayment{}, err
	}
	rerr := copier.Copy(&res.Payments, &payments)
	if rerr != nil {
		return api.ResGetPayment{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	return res, nil
}

func getPaymentList(dbClient db.Client, offset *int, year *int) ([]paymtent, *response.Error) {
	searchYear := utils.GetSchoolYear()
	if year != nil {
		searchYear = *year
	}
	params := struct {
		Offset *int `twowaysql:"offset"`
		Year   int  `twowaysql:"year"`
	}{
		Offset: offset,
		Year:   searchYear,
	}
	paymtents := []paymtent{}
	err := dbClient.Select(&paymtents, "sql/payment/select_payment.sql", &params)
	if err != nil {
		return []paymtent{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(paymtents) == 0 {
		return []paymtent{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "支払いが存在しません", Log: "no rows in result"}
	}
	return paymtents, nil
}
