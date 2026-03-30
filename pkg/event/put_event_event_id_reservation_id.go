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
	// 日付のビジネスロジック検証（共通関数を使用）
	if err := validateReservationDates(requestBody.StartDate, requestBody.FinishDate, requestBody.ReservationStartDate, requestBody.ReservationFinishDate); err != nil {
		return err
	}

	// 新しいcapacityが既に予約済み人数以上であることを検証
	currentReservationCount, err := getReservationCountFromReservationId(dbClient, eventId, reservationId)
	if err != nil {
		return err
	}
	if requestBody.Capacity < currentReservationCount {
		return &response.Error{Code: http.StatusConflict, Level: "Info", Message: "参加枠は既に予約済みの人数以上である必要があります", Log: "requested capacity is less than current reservation count"}
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
	rowsAffected, rowErr := result.RowsAffected()
	if rowErr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "更新結果の確認に失敗しました", Log: rowErr.Error()}
	}
	if rowsAffected == 0 {
		return &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "指定された予約枠が見つかりません", Log: "reservation not found"}
	}

	return nil
}

// [must] 指定された予約枠の現在の予約数を取得する
func getReservationCountFromReservationId(dbClient db.TransactionClient, eventId string, reservationId string) (int, *response.Error) {
	params := struct {
		EventId       string `twowaysql:"eventId"`
		ReservationId string `twowaysql:"reservationId"`
	}{
		EventId:       eventId,
		ReservationId: reservationId,
	}

	type reservationCount struct {
		Count int `db:"count"`
	}

	results := []reservationCount{}
	err := dbClient.Select(&results, "sql/event/select_reservation_count_from_reservation_id.sql", &params)
	if err != nil {
		return 0, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約数の取得に失敗しました", Log: err.Error()}
	}

	if len(results) == 0 {
		return 0, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "指定された予約枠が見つかりません", Log: "reservation not found"}
	}

	return results[0].Count, nil
}
