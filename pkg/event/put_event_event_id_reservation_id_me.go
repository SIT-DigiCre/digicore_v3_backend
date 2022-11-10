package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutEventEventIDReservationIDMe(ctx echo.Context, dbClient db.TransactionClient, eventID string, reservationID string, requestBody api.ReqPutEventEventIDReservationIDMe) (api.ResGetEventEventIDReservationID, *response.Error) {
	userID := ctx.Get("user_id").(string)
	err := updateReservationUser(dbClient, reservationID, userID, requestBody)
	if err != nil {
		return api.ResGetEventEventIDReservationID{}, err
	}

	return GetEventEventIDReservationID(ctx, dbClient, eventID, reservationID)
}

func updateReservationUser(dbClient db.TransactionClient, reservationID string, userID string, requestBody api.ReqPutEventEventIDReservationIDMe) *response.Error {
	params := struct {
		ReservationID string `twowaysql:"reservationID"`
		UserID        string `twowaysql:"userID"`
		URL           string `twowaysql:"URL"`
		Comment       string `twowaysql:"comment"`
	}{
		ReservationID: reservationID,
		UserID:        userID,
		URL:           requestBody.URL,
		Comment:       requestBody.Comment,
	}
	_, err := dbClient.DuplicateUpdate("sql/event/insert_reservation_user.sql", "sql/event/update_reservation_user.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
