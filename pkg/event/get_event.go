package event

import (
	"fmt"
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
	event, err := EventFromUserID(dbClient, userID)
	if err != nil {
		return api.ResGetEvent{}, err
	}
	fmt.Println(event)
	rerr := copier.Copy(&res.Event, &event)
	if rerr != nil {
		return api.ResGetEvent{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "プロフィールの読み込みに失敗しました", Log: rerr.Error()}
	}
	return res, nil
}

type event struct {
	EventID    string `db:"event_id"`
	Name       string `db:"name"`
	Reservable bool   `db:"reservable"`
	Reservated bool   `db:"reservated"`
}

func EventFromUserID(dbClient db.Client, userID string) ([]event, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	event := []event{}
	err := dbClient.Select(&event, "sql/event/select_event_list.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "イベント一覧の取得に失敗しました", Log: err.Error()}
	}
	return event, nil
}
