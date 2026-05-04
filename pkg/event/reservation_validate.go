package event

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
)

// validateReservationDates は予約枠の日付ビジネスロジックを検証する共通関数。
// Start < Finish, ReservationStart < ReservationFinish, ReservationStart <= Finish を確認。
func validateReservationDates(startDate, finishDate, reservationStartDate, reservationFinishDate time.Time) *response.Error {
	// 開始日時は終了日時より前である必要がある
	if !startDate.Before(finishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "開始日時は終了日時より前である必要があります", Log: "startDate is not before finishDate"}
	}

	// 予約開始日時は予約終了日時より前である必要がある
	if !reservationStartDate.Before(reservationFinishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約開始日時は予約終了日時より前である必要があります", Log: "reservationStartDate is not before reservationFinishDate"}
	}

	// 予約開始日時はイベント終了日時以前である必要がある
	if reservationStartDate.After(finishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約開始日時はイベント終了日時以前である必要があります", Log: "reservationStartDate is after finishDate"}
	}

	return nil
}
