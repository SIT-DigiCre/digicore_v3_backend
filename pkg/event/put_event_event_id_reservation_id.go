package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// 予約枠の情報を更新する関数。
// 更新に成功した場合、更新後の予約枠の詳細を返す。
func PutEventEventIdReservationId(ctx echo.Context, dbClient db.TransactionClient, eventId string, reservationId string, requestBody api.PutEventEventIdReservationIdJSONRequestBody) (api.ResGetEventEventIdReservationId, *response.Error) {
	// 予約枠を更新
	if err := updateEventReservation(dbClient, eventId, reservationId, requestBody); err != nil {
		return api.ResGetEventEventIdReservationId{}, err
	}

	// UUIDのパース
	eventUUID, parseErr := uuid.Parse(eventId)
	if parseErr != nil {
		return api.ResGetEventEventIdReservationId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "イベントIDの解析に失敗しました", Log: parseErr.Error()}
	}

	reservationUUID, parseErr := uuid.Parse(reservationId)
	if parseErr != nil {
		return api.ResGetEventEventIdReservationId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠IDの解析に失敗しました", Log: parseErr.Error()}
	}

	// 更新した予約枠情報を返す
	return api.ResGetEventEventIdReservationId{
		EventId:               eventUUID.String(),
		ReservationId:         reservationUUID.String(),
		Name:                  requestBody.Name,
		Description:           requestBody.Description,
		StartDate:             requestBody.StartDate.Format("2006-01-02T15:04:05Z07:00"),
		FinishDate:            requestBody.FinishDate.Format("2006-01-02T15:04:05Z07:00"),
		ReservationStartDate:  requestBody.ReservationStartDate.Format("2006-01-02T15:04:05Z07:00"),
		ReservationFinishDate: requestBody.ReservationFinishDate.Format("2006-01-02T15:04:05Z07:00"),
		Capacity:              requestBody.Capacity,
		Reservable:            false,
		Reservated:            false,
		FreeCapacity:          requestBody.Capacity,
		Users:                 []api.ResGetEventEventIdReservationIdObjectUser{},
	}, nil
}

func updateEventReservation(dbClient db.TransactionClient, eventId string, reservationId string, requestBody api.PutEventEventIdReservationIdJSONRequestBody) *response.Error {
	params := struct {
		EventId               string `twowaysql:"eventId"`
		ReservationId         string `twowaysql:"reservationId"`
		Name                  string `twowaysql:"name"`
		Description           string `twowaysql:"description"`
		StartDate             string `twowaysql:"startDate"`
		FinishDate            string `twowaysql:"finishDate"`
		ReservationStartDate  string `twowaysql:"reservationStartDate"`
		ReservationFinishDate string `twowaysql:"reservationFinishDate"`
		Capacity              int    `twowaysql:"capacity"`
	}{
		EventId:               eventId,
		ReservationId:         reservationId,
		Name:                  requestBody.Name,
		Description:           requestBody.Description,
		StartDate:             requestBody.StartDate.Format("2006-01-02 15:04:05"),
		FinishDate:            requestBody.FinishDate.Format("2006-01-02 15:04:05"),
		ReservationStartDate:  requestBody.ReservationStartDate.Format("2006-01-02 15:04:05"),
		ReservationFinishDate: requestBody.ReservationFinishDate.Format("2006-01-02 15:04:05"),
		Capacity:              requestBody.Capacity,
	}
	result, rerr := dbClient.Exec("sql/event/update_reservation.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}

	// 影響を受けた行数を確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "更新結果の確認に失敗しました", Log: err.Error()}
	}
	if rowsAffected == 0 {
		return &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "指定された予約枠が見つかりません", Log: "reservation not found"}
	}

	return nil
}
