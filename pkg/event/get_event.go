package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetEvent(ctx echo.Context, dbClient db.Client, params api.GetEventParams) (api.ResGetEvent, *response.Error) {
	res := api.ResGetEvent{}
	userID := ctx.Get("user_id").(string)
	events, err := eventList(dbClient, userID, params.Offset)
	if err != nil {
		return api.ResGetEvent{}, err
	}
	rerr := copier.Copy(&res.Event, &events)
	if rerr != nil {
		return api.ResGetEvent{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベント一覧の取得に失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type event struct {
	EventID      string `db:"event_id"`
	Name         string `db:"name"`
	CalendarView bool   `db:"calendar_view"`
	Reservable   bool   `db:"reservable"`
	Reservated   bool   `db:"reservated"`
}

func eventList(dbClient db.Client, userID string, offset *int) ([]event, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
		Offset *int   `twowaysql:"offset"`
	}{
		UserID: userID,
		Offset: offset,
	}
	events := []event{}
	err := dbClient.Select(&events, "sql/event/select_event_list.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベント一覧の取得に失敗しました", Log: err.Error()}
	}
	if len(events) == 0 {
		return nil, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "イベントがありません。", Log: "no rows in result"}
	}
	return events, nil
}
