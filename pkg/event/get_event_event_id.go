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

func GetEventEventID(ctx echo.Context, dbClient db.Client, eventID string) (api.ResGetEventEventID, *response.Error) {
	res := api.ResGetEventEventID{}
	userID := ctx.Get("user_id").(string)
	eventDetail, err := getEventFromEventID(dbClient, eventID, userID)
	if err != nil {
		return api.ResGetEventEventID{}, err
	}
	rerr := copier.Copy(&res, &eventDetail)
	if rerr != nil {
		return api.ResGetEventEventID{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベント一覧の取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type eventDetail struct {
	EventID      string `db:"event_id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	CalendarView bool   `db:"calendar_view"`
	Reservable   bool   `db:"reservable"`
	Reservated   bool   `db:"reservated"`
	Reservation  []eventDetailObjectReservation
}

type eventDetailObjectReservation struct {
	Capacity              int    `db:"capacity"`
	Description           string `db:"description"`
	FreeCapacity          int    `db:"free_capacity"`
	Name                  string `db:"name"`
	Reservable            bool   `db:"reservable"`
	Reservated            bool   `db:"reservated"`
	ReservationID         string `db:"reservation_id"`
	StartDate             string `db:"start_date"`
	FinishDate            string `db:"finish_date"`
	ReservationStartDate  string `db:"reservation_start_date"`
	ReservationFinishDate string `db:"reservation_finish_date"`
}

func getEventFromEventID(dbClient db.Client, eventID string, userID string) (eventDetail, *response.Error) {
	params := struct {
		EventID string `twowaysql:"eventID"`
		UserID  string `twowaysql:"userID"`
	}{
		EventID: eventID,
		UserID:  userID,
	}
	eventDetails := []eventDetail{}
	err := dbClient.Select(&eventDetails, "sql/event/select_event_from_event_id.sql", &params)
	if err != nil {
		return eventDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	if len(eventDetails) == 0 {
		return eventDetail{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "イベントがありません。", Log: "no rows in result"}
	}
	eventReservations := []eventDetailObjectReservation{}
	err = dbClient.Select(&eventReservations, "sql/event/select_event_reservation_from_event_id.sql", &params)
	if err != nil {
		return eventDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
	}
	now_date := time.Now()
	for i := range eventReservations {
		reservationStartDate, err := time.Parse(time.RFC3339, eventReservations[i].ReservationStartDate)
		if err != nil {
			return eventDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
		}
		reservationFinishDate, err := time.Parse(time.RFC3339, eventReservations[i].ReservationFinishDate)
		if err != nil {
			return eventDetail{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの取得に失敗しました", Log: err.Error()}
		}
		if eventReservations[i].FreeCapacity == 0 || now_date.Before(reservationStartDate) || now_date.After(reservationFinishDate) {
			eventReservations[i].Reservable = false
		} else {
			eventReservations[i].Reservable = true
		}
	}
	eventDetails[0].Reservation = eventReservations
	return eventDetails[0], nil
}
