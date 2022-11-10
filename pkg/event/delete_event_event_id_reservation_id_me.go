package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteEventEventIDReservationIDMe(ctx echo.Context, dbClient db.TransactionClient, eventID string, reservationID string) (api.ResGetEventEventIDReservationID, *response.Error) {
	userID := ctx.Get("user_id").(string)
	err := deleteReservationUser(dbClient, reservationID, userID)
	if err != nil {
		return api.ResGetEventEventIDReservationID{}, err
	}

	return GetEventEventIDReservationID(ctx, dbClient, eventID, reservationID)
}

func deleteReservationUser(dbClient db.TransactionClient, reservationID string, userID string) *response.Error {
	params := struct {
		ReservationID string `twowaysql:"reservationID"`
		UserID        string `twowaysql:"userID"`
	}{
		ReservationID: reservationID,
		UserID:        userID,
	}
	_, err := dbClient.Exec("sql/event/delete_reservation_user.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
