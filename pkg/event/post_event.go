package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// name, description, calendar_viewを受け取ってイベントを作成する関数。
// 作成に成功した場合、イベントの詳細を返す。
func PostEvent(ctx echo.Context, dbClient db.TransactionClient, requestBody api.PostEventJSONRequestBody) (api.ResPostEventEvent, *response.Error) {
	// イベントを作成し、イベントIDを取得
	eventId, err := createEvent(dbClient, requestBody)
	if err != nil {
		return api.ResPostEventEvent{}, err
	}

	// UUIDのパース
	eventUUID, parseErr := uuid.Parse(eventId)
	if parseErr != nil {
		return api.ResPostEventEvent{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントIDの解析に失敗しました", Log: parseErr.Error()}
	}

	// 作成したイベント情報を返す
	return api.ResPostEventEvent{
		EventId:      eventUUID,
		Name:         requestBody.Name,
		Description:  requestBody.Description,
		CalendarView: requestBody.CalendarView,
	}, nil
}

func createEvent(dbClient db.TransactionClient, requestBody api.PostEventJSONRequestBody) (string, *response.Error) {
	// イベントIDを生成
	_, rerr := dbClient.Exec("sql/transaction/generate_id.sql", nil, false)
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの作成に失敗しました", Log: rerr.Error()}
	}
	eventId, gerr := dbClient.GetId()
	if gerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントの作成に失敗しました", Log: gerr.Error()}
	}
	if eventId == "" {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントIDの生成に失敗しました", Log: "空のIDが返されました"}
	}

	// eventsテーブルに挿入
	eventParams := struct {
		EventId      string `twowaysql:"eventId"`
		Name         string `twowaysql:"name"`
		Description  string `twowaysql:"description"`
		CalendarView bool   `twowaysql:"calendarView"`
	}{
		EventId:      eventId,
		Name:         requestBody.Name,
		Description:  requestBody.Description,
		CalendarView: requestBody.CalendarView,
	}
	logrus.Debugf("Inserting event: sql=insert_event.sql params=%#v", eventParams)
	_, rerr = dbClient.Exec("sql/event/insert_event.sql", &eventParams, false)
	if rerr != nil {
		logrus.Errorf("イベントの挿入に失敗しました: %v", rerr)
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました (events)", Log: rerr.Error()}
	}

	return eventId, nil
}
