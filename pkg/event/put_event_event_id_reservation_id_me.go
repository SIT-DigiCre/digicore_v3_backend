package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func PutEventEventIdReservationIdMe(ctx echo.Context, dbClient db.TransactionClient, eventId string, reservationId string, requestBody api.ReqPutEventEventIdReservationIdMe) (api.ResGetEventEventIdReservationId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	err := updateReservationUser(dbClient, reservationId, userId, requestBody)
	if err != nil {
		return api.ResGetEventEventIdReservationId{}, err
	}

	return GetEventEventIdReservationId(ctx, dbClient, eventId, reservationId)
}

func updateReservationUser(dbClient db.TransactionClient, reservationId string, userId string, requestBody api.ReqPutEventEventIdReservationIdMe) *response.Error {
	params := struct {
		ReservationId string `twowaysql:"reservationId"`
		UserId        string `twowaysql:"userId"`
		Url           string `twowaysql:"url"`
		Comment       string `twowaysql:"comment"`
	}{
		ReservationId: reservationId,
		UserId:        userId,
		Url:           requestBody.Url,
		Comment:       requestBody.Comment,
	}
	_, _, err := dbClient.DuplicateUpdate("sql/event/insert_reservation_user.sql", "sql/event/update_reservation_user.sql", &params)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}
