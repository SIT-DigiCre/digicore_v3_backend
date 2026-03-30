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

// レスポンス型
type ResPostEventEventIdReservation struct {
	EventId               uuid.UUID `json:"eventId"`
	ReservationId         uuid.UUID `json:"reservationId"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	StartDate             string    `json:"startDate"`
	FinishDate            string    `json:"finishDate"`
	ReservationStartDate  string    `json:"reservationStartDate"`
	ReservationFinishDate string    `json:"reservationFinishDate"`
	Capacity              int       `json:"capacity"`
}

// eventId 配下に新しい予約枠を作成する関数。
// 作成に成功した場合、作成した予約枠の詳細を返す。
func PostEventEventIdReservation(ctx echo.Context, dbClient db.TransactionClient, eventId string, requestBody api.PostEventEventIdReservationJSONRequestBody) (ResPostEventEventIdReservation, *response.Error) {
	// 予約枠を作成し、予約枠IDを取得
	reservationId, err := createEventReservation(dbClient, eventId, requestBody)
	if err != nil {
		return ResPostEventEventIdReservation{}, err
	}

	// UUIDのパース
	eventUUID, parseErr := uuid.Parse(eventId)
	if parseErr != nil {
		return ResPostEventEventIdReservation{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントIDの解析に失敗しました", Log: parseErr.Error()}
	}

	reservationUUID, parseErr := uuid.Parse(reservationId)
	if parseErr != nil {
		return ResPostEventEventIdReservation{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠IDの解析に失敗しました", Log: parseErr.Error()}
	}

	// 作成した予約枠情報を返す
	return ResPostEventEventIdReservation{
		EventId:               eventUUID,
		ReservationId:         reservationUUID,
		Name:                  requestBody.Name,
		Description:           requestBody.Description,
		StartDate:             requestBody.StartDate.Format("2006-01-02T15:04:05Z07:00"),
		FinishDate:            requestBody.FinishDate.Format("2006-01-02T15:04:05Z07:00"),
		ReservationStartDate:  requestBody.ReservationStartDate.Format("2006-01-02T15:04:05Z07:00"),
		ReservationFinishDate: requestBody.ReservationFinishDate.Format("2006-01-02T15:04:05Z07:00"),
		Capacity:              requestBody.Capacity,
	}, nil
}

func createEventReservation(dbClient db.TransactionClient, eventId string, requestBody api.PostEventEventIdReservationJSONRequestBody) (string, *response.Error) {
	// 日付のビジネスロジック検証
	if !requestBody.StartDate.Before(requestBody.FinishDate) {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "開始日時は終了日時より前である必要があります", Log: "startDate is not before finishDate"}
	}
	if !requestBody.ReservationStartDate.Before(requestBody.ReservationFinishDate) {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約開始日時は予約終了日時より前である必要があります", Log: "reservationStartDate is not before reservationFinishDate"}
	}
	if requestBody.ReservationStartDate.After(requestBody.FinishDate) {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約開始日時はイベント終了日時以前である必要があります", Log: "reservationStartDate is after finishDate"}
	}

	// 予約枠IDを生成
	_, rerr := dbClient.Exec("sql/transaction/generate_id.sql", nil, false)
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠の作成に失敗しました", Log: rerr.Error()}
	}
	reservationId, gerr := dbClient.GetId()
	if gerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠の作成に失敗しました", Log: gerr.Error()}
	}
	if reservationId == "" {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠IDの生成に失敗しました", Log: "空のIDが返されました"}
	}

	// event_reservationsテーブルに挿入
	reservationParams := struct {
		ReservationId         string `twowaysql:"reservationId"`
		EventId               string `twowaysql:"eventId"`
		Name                  string `twowaysql:"name"`
		Description           string `twowaysql:"description"`
		StartDate             string `twowaysql:"startDate"`
		FinishDate            string `twowaysql:"finishDate"`
		ReservationStartDate  string `twowaysql:"reservationStartDate"`
		ReservationFinishDate string `twowaysql:"reservationFinishDate"`
		Capacity              int    `twowaysql:"capacity"`
	}{
		ReservationId:         reservationId,
		EventId:               eventId,
		Name:                  requestBody.Name,
		Description:           requestBody.Description,
		StartDate:             requestBody.StartDate.Format("2006-01-02 15:04:05"),
		FinishDate:            requestBody.FinishDate.Format("2006-01-02 15:04:05"),
		ReservationStartDate:  requestBody.ReservationStartDate.Format("2006-01-02 15:04:05"),
		ReservationFinishDate: requestBody.ReservationFinishDate.Format("2006-01-02 15:04:05"),
		Capacity:              requestBody.Capacity,
	}
	logrus.Debugf("Inserting event reservation: sql=insert_reservation.sql params=%#v", reservationParams)
	_, rerr = dbClient.Exec("sql/event/insert_reservation.sql", &reservationParams, false)
	if rerr != nil {
		logrus.Errorf("予約枠の挿入に失敗しました: %v", rerr)
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました (event_reservations)", Log: rerr.Error()}
	}

	return reservationId, nil
}
