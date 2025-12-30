package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// イベントのデータを更新する関数。
// 更新に成功した場合、更新後のイベントの詳細を返す。
func PutEvent(ctx echo.Context, dbClient db.TransactionClient, eventId string, requestBody api.PutEventEventIdJSONRequestBody) (api.ResPostEventEvent, *response.Error) {
	// イベントのメタデータを更新
	if err := updateEvent(dbClient, eventId, requestBody); err != nil {
		return api.ResPostEventEvent{}, err
	}

	// UUIDのパース
	eventUUID, parseErr := uuid.Parse(eventId)
	if parseErr != nil {
		return api.ResPostEventEvent{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントIDの解析に失敗しました", Log: parseErr.Error()}
	}

	// 更新したイベント情報を返す
	return api.ResPostEventEvent{
		EventId:      eventUUID,
		Name:         requestBody.Name,
		Description:  requestBody.Description,
		CalendarView: requestBody.CalendarView,
	}, nil
}

func updateEvent(dbClient db.TransactionClient, eventId string, requestBody api.PutEventEventIdJSONRequestBody) *response.Error {
	params := struct {
		EventId      string `twowaysql:"eventId"`
		Name         string `twowaysql:"name"`
		Description  string `twowaysql:"description"`
		CalendarView int    `twowaysql:"calendarView"`
	}{
		EventId:      eventId,
		Name:         requestBody.Name,
		Description:  requestBody.Description,
		CalendarView: boolToInt(requestBody.CalendarView),
	}
	_, rerr := dbClient.Exec("sql/event/update_event.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	return nil
}
