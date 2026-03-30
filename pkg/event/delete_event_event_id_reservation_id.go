package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

// 予約枠を削除する関数。
// 削除に成功した場合、エラーを返さない。
func DeleteEventEventIdReservationId(ctx echo.Context, dbClient db.TransactionClient, eventId string, reservationId string) *response.Error {
	// 参加者を削除
	deleteUsersParams := struct {
		ReservationId string `twowaysql:"reservationId"`
	}{
		ReservationId: reservationId,
	}
	_, uerr := dbClient.Exec("sql/event/delete_reservation_users.sql", &deleteUsersParams, false)
	if uerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました", Log: uerr.Error()}
	}

	// 予約枠を削除
	params := struct {
		EventId       string `twowaysql:"eventId"`
		ReservationId string `twowaysql:"reservationId"`
	}{
		EventId:       eventId,
		ReservationId: reservationId,
	}
	result, rerr := dbClient.Exec("sql/event/delete_reservation.sql", &params, false)
	if rerr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}

	// 影響を受けた行数を確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "削除結果の確認に失敗しました", Log: err.Error()}
	}
	if rowsAffected == 0 {
		return &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "指定された予約枠が見つかりません", Log: "reservation not found"}
	}

	return nil
}
