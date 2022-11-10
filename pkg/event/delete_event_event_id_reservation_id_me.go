package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func DeleteEventEventIdReservationIdMe(ctx echo.Context, dbClient db.TransactionClient, eventId string, reservationId string) (api.ResGetEventEventIdReservationId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := deleteReservationUser(dbClient, reservationId, userId)
	if err != nil {
		return api.ResGetEventEventIdReservationId{}, err
	}

	return GetEventEventIdReservationId(ctx, dbClient, eventId, reservationId)
}

func deleteReservationUser(dbClient db.TransactionClient, reservationId string, userId string) *response.Error {
	params := struct {
		ReservationId string `twowaysql:"reservationId"`
		UserId        string `twowaysql:"userId"`
	}{
		ReservationId: reservationId,
		UserId:        userId,
	}
	_, err := dbClient.Exec("sql/event/delete_reservation_user.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
