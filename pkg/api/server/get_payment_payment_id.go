package server

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/payment"
	"github.com/labstack/echo/v4"
)

func (s *server) GetPaymentPaymentId(ctx echo.Context, paymentId string) error {
	dbClient := db.Open()

	res, err := payment.GetPaymentPaymentId(ctx, &dbClient, paymentId)
	if err != nil {
		return response.ErrorResponse(ctx, err)
	}

	return response.SuccessResponse(ctx, res)
}
