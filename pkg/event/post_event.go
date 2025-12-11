package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// PostEvent creates an event and a single reservation attached to it.
// It returns the full event detail via GetEventEventId on success.
func PostEvent(ctx echo.Context, dbClient db.TransactionClient, requestBody api.PostEventJSONBody) (api.ResGetEventEventId, *response.Error) {
	// create event and reservation and get event id
	eventId, err := createEvent(dbClient, requestBody)
	if err != nil {
		return api.ResGetEventEventId{}, err
	}
	return GetEventEventId(ctx, dbClient, eventId)
}

func createEvent(dbClient db.TransactionClient, requestBody api.PostEventJSONBody) (string, *response.Error) {
	// generate event id
	_, rerr := dbClient.Exec("sql/transaction/generate_id.sql", nil, false)
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの作成に失敗しました", Log: rerr.Error()}
	}
	eventId, gerr := dbClient.GetId()
	if gerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの作成に失敗しました", Log: gerr.Error()}
	}
	if eventId == "" {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントIDの生成に失敗しました", Log: "empty id returned"}
	}

	// generate reservation id
	_, rerr = dbClient.Exec("sql/transaction/generate_id.sql", nil, false)
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの作成に失敗しました", Log: rerr.Error()}
	}
	reservationId, gerr := dbClient.GetId()
	if gerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの作成に失敗しました", Log: gerr.Error()}
	}
	if reservationId == "" {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約IDの生成に失敗しました", Log: "empty id returned"}
	}

	// insert into events table (use struct with twowaysql tags)
	eventParams := struct {
		EventId     string `twowaysql:"eventId"`
		Name        string `twowaysql:"name"`
		Description string `twowaysql:"description"`
	}{
		EventId:     eventId,
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}
	logrus.Debugf("Inserting event: sql=insert_events.sql params=%#v", eventParams)
	_, rerr = dbClient.Exec("sql/event/insert_events.sql", &eventParams, false)
	if rerr != nil {
		logrus.Errorf("insert_events failed: %v", rerr)
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました (events)", Log: rerr.Error()}
	}

	// insert into event_reservations table
	reservationParams := struct {
		ReservationId     string `twowaysql:"reservationId"`
		EventId           string `twowaysql:"eventId"`
		Name              string `twowaysql:"name"`
		Description       string `twowaysql:"description"`
		StartDate         string `twowaysql:"startDate"`
		FinishDate        string `twowaysql:"finishDate"`
		ReservationStart  string `twowaysql:"reservationStart"`
		ReservationFinish string `twowaysql:"reservationFinish"`
		Capacity          int    `twowaysql:"capacity"`
	}{
		ReservationId:     reservationId,
		EventId:           eventId,
		Name:              requestBody.Name,
		Description:       requestBody.Description,
		StartDate:         requestBody.StartDate.Format("2006-01-02 15:04:05"),
		FinishDate:        requestBody.FinishDate.Format("2006-01-02 15:04:05"),
		ReservationStart:  requestBody.ReservationStart.Format("2006-01-02 15:04:05"),
		ReservationFinish: requestBody.ReservationFinish.Format("2006-01-02 15:04:05"),
		Capacity:          requestBody.Capacity,
	}
	logrus.Debugf("Inserting reservation: sql=insert_event_reservations.sql params=%#v", reservationParams)
	_, rerr = dbClient.Exec("sql/event/insert_event_reservations.sql", &reservationParams, false)
	if rerr != nil {
		logrus.Errorf("insert_event_reservations failed: %v", rerr)
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました (reservations)", Log: rerr.Error()}
	}
	return eventId, nil
}
