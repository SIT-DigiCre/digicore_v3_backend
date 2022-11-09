package event

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetEventEventIDReservationID(ctx echo.Context, dbClient db.Client, eventID string, reservationID string) (api.ResGetEventEventIDReservationID, *response.Error) {
	res := api.ResGetEventEventIDReservationID{}
	userID := ctx.Get("user_id").(string)
	reservation, err := getReservationFromReservationID(dbClient, eventID, reservationID, userID)
	if err != nil {
		return api.ResGetEventEventIDReservationID{}, err
	}
	rerr := copier.Copy(&res, &reservation)
	if rerr != nil {
		return api.ResGetEventEventIDReservationID{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの予約枠の取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type eventReservation struct {
	EventID               string `db:"event_id"`
	ReservationID         string `db:"reservation_id"`
	Name                  string `db:"name"`
	Reservable            bool   `db:"reservable"`
	Reservated            bool   `db:"reservated"`
	Capacity              int    `db:"capacity"`
	Description           string `db:"description"`
	FreeCapacity          int    `db:"free_capacity"`
	StartDate             string `db:"start_date"`
	FinishDate            string `db:"finish_date"`
	ReservationStartDate  string `db:"reservation_start_date"`
	ReservationFinishDate string `db:"reservation_finish_date"`
	User                  []eventReservationObjectUser
}

type eventReservationObjectUser struct {
	URL      string `db:"url"`
	Comment  string `db:"comment"`
	Name     string `db:"username"`
	UserIcon string `db:"icon_url"`
	UserID   string `db:"user_id"`
}

func getReservationFromReservationID(dbClient db.Client, eventID string, reservationID string, userID string) (eventReservation, *response.Error) {
	params := struct {
		EventID       string `twowaysql:"eventID"`
		ReservationID string `twowaysql:"reservationID"`
		UserID        string `twowaysql:"userID"`
	}{
		EventID:       eventID,
		ReservationID: reservationID,
		UserID:        userID,
	}
	eventReservations := []eventReservation{}
	err := dbClient.Select(&eventReservations, "sql/event/select_event_reservation_from_event_id_reservation_id.sql", &params)
	if err != nil {
		return eventReservation{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	now_date := time.Now()
	reservationStartDate, err := time.Parse(time.RFC3339, eventReservations[0].ReservationStartDate)
	if err != nil {
		return eventReservation{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	reservationFinishDate, err := time.Parse(time.RFC3339, eventReservations[0].ReservationFinishDate)
	if err != nil {
		return eventReservation{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	if eventReservations[0].FreeCapacity == 0 || now_date.Before(reservationStartDate) || now_date.After(reservationFinishDate) {
		eventReservations[0].Reservable = false
	} else {
		eventReservations[0].Reservable = true
	}
	eventReservationObjectUsers := []eventReservationObjectUser{}
	err = dbClient.Select(&eventReservationObjectUsers, "sql/event/select_event_reservation_user_from_event_id_reservation_id.sql", &params)
	if err != nil {
		return eventReservation{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	eventReservations[0].User = eventReservationObjectUsers
	return eventReservations[0], nil
}
