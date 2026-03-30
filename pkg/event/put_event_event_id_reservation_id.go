package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

// 予約枠の情報を更新する関数。
// 更新に成功した場合、更新後の予約枠の詳細を返す。
func PutEventEventIdReservationId(ctx echo.Context, dbClient db.TransactionClient, eventId string, reservationId string, requestBody api.PutEventEventIdReservationIdJSONRequestBody) (api.ResGetEventEventIdReservationId, *response.Error) {
	// 予約枠を更新
	if err := updateEventReservation(dbClient, eventId, reservationId, requestBody); err != nil {
		return api.ResGetEventEventIdReservationId{}, err
	}

	// 更新後の予約枠情報をデータベースから取得して返す
	reservation, err := getReservationFromReservationId(dbClient, eventId, reservationId, ctx.Get("user_id").(string))
	if err != nil {
		return api.ResGetEventEventIdReservationId{}, err
	}

	// eventReservationをapi.ResGetEventEventIdReservationIdに変換
	res := api.ResGetEventEventIdReservationId{
		EventId:               reservation.EventId,
		ReservationId:         reservation.ReservationId,
		Name:                  reservation.Name,
		Description:           reservation.Description,
		StartDate:             reservation.StartDate,
		FinishDate:            reservation.FinishDate,
		ReservationStartDate:  reservation.ReservationStartDate,
		ReservationFinishDate: reservation.ReservationFinishDate,
		Capacity:              reservation.Capacity,
		Reservable:            reservation.Reservable,
		Reservated:            reservation.Reservated,
		FreeCapacity:          reservation.FreeCapacity,
	}

	// ユーザー情報の変換
	if reservation.Users != nil {
		for _, user := range reservation.Users {
			res.Users = append(res.Users, api.ResGetEventEventIdReservationIdObjectUser{
				Url:       user.Url,
				Comment:   user.Comment,
				Name:      user.Name,
				UserIcon:  user.UserIcon,
				UserId:    user.UserId,
				CreatedAt: user.CreatedAt,
			})
		}
	} else {
		res.Users = []api.ResGetEventEventIdReservationIdObjectUser{}
	}

	return res, nil
}

func updateEventReservation(dbClient db.TransactionClient, eventId string, reservationId string, requestBody api.PutEventEventIdReservationIdJSONRequestBody) *response.Error {
	// 日付のビジネスロジック検証
	if !requestBody.StartDate.Before(requestBody.FinishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "開始日時は終了日時より前である必要があります", Log: "startDate is not before finishDate"}
	}
	if !requestBody.ReservationStartDate.Before(requestBody.ReservationFinishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約開始日時は予約終了日時より前である必要があります", Log: "reservationStartDate is not before reservationFinishDate"}
	}
	if requestBody.ReservationStartDate.After(requestBody.FinishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約開始日時はイベント終了日時以前である必要があります", Log: "reservationStartDate is after finishDate"}
	}

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
