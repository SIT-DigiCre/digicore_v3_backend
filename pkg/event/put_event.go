package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

// PutEvent updates event metadata and the (single) reservation for the event.
// It returns the full event detail via GetEventEventId on success.
func PutEvent(ctx echo.Context, dbClient db.TransactionClient, eventId string, requestBody api.PostEventJSONBody) (api.ResGetEventEventId, *response.Error) {
	// update event meta
	if err := updateEvent(dbClient, eventId, requestBody); err != nil {
		return api.ResGetEventEventId{}, err
	}
	// update reservation
	if err := updateEventReservation(dbClient, eventId, requestBody); err != nil {
		return api.ResGetEventEventId{}, err
	}
	return GetEventEventId(ctx, dbClient, eventId)
}

func updateEvent(dbClient db.TransactionClient, eventId string, requestBody api.PostEventJSONBody) *response.Error {
	params := struct {
		EventId     string `twowaysql:"eventId"`
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		EventId:     eventId,
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	_, rerr := dbClient.Exec("sql/event/update_event.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}

func updateEventReservation(dbClient db.TransactionClient, eventId string, requestBody api.PostEventJSONBody) *response.Error {
	// get reservation id for this event
	params := struct {
		EventId string `twowaysql:"eventId"`
	}{EventId: eventId}

	reservations := []struct {
		ReservationId string `db:"reservation_id"`
	}{}
	if err := dbClient.Select(&reservations, "sql/event/select_event_reservation_from_event_id.sql", &params); err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	if len(reservations) == 0 {
		return &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "予約が見つかりません", Log: "no rows in result"}
	}
	reservationId := reservations[0].ReservationId

	up := struct {
		ReservationId     string `twowaysql:"reservationId"`
		Name              string `twowaysql:"name"`
		Description       string `twowaysql:"description"`
		StartDate         string `twowaysql:"startDate"`
		FinishDate        string `twowaysql:"finishDate"`
		ReservationStart  string `twowaysql:"reservationStart"`
		ReservationFinish string `twowaysql:"reservationFinish"`
		Capacity          int    `twowaysql:"capacity"`
	}{
		ReservationId:     reservationId,
		Name:              requestBody.Name,
		Description:       requestBody.Description,
	StartDate:         requestBody.StartDate.Format("2006-01-02 15:04:05"),
	FinishDate:        requestBody.FinishDate.Format("2006-01-02 15:04:05"),
	ReservationStart:  requestBody.ReservationStart.Format("2006-01-02 15:04:05"),
	ReservationFinish: requestBody.ReservationFinish.Format("2006-01-02 15:04:05"),
		Capacity:          requestBody.Capacity,
	}

	_, rerr := dbClient.Exec("sql/event/update_event_reservation.sql", &up, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
