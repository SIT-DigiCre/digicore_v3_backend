package payment

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func PutPaymentPaymentId(ctx echo.Context, dbClient db.TransactionClient, paymentId string, requestBody api.ReqPutPaymentPaymentId) (api.ResGetPaymentPaymentId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := updatePayment(dbClient, paymentId, requestBody)
	if err != nil {
		return api.ResGetPaymentPaymentId{}, err
	}
	err = utils.RenewalActiveLimit(dbClient, userId, strconv.Itoa(utils.GetSchoolYear()+1)+"-05-01")
	if err != nil {
		return api.ResGetPaymentPaymentId{}, err
	}
	utils.NoticeMattermost(fmt.Sprintf("部費振込申請(%s)が行われました", userId), "digicore-notice", "digicore-notice", "bell")
	return GetPaymentPaymentId(ctx, dbClient, paymentId)
}

func updatePayment(dbClient db.TransactionClient, paymentId string, requestBody api.ReqPutPaymentPaymentId) *response.Error {
	params := struct {
		PaymentId string `twowaysql:"paymentId"`
		Checked   bool   `twowaysql:"checked"`
		Note      string `twowaysql:"note"`
	}{
		PaymentId: paymentId,
		Checked:   requestBody.Checked,
		Note:      requestBody.Note,
	}
	_, err := dbClient.Exec("sql/payment/update_payment.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
