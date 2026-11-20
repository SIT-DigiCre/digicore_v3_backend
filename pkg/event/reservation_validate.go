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

// validateReservationPeriod は予約の申込・キャンセルが可能な期間かを検証する共通関数。
// 予約期間内かつ予約枠の終了前である必要がある。
func validateReservationPeriod(reservation eventReservation) *response.Error {
	reservationStartDate, terr := time.Parse(time.RFC3339, reservation.ReservationStartDate)
	if terr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠の取得に失敗しました", Log: terr.Error()}
	}
	reservationFinishDate, terr := time.Parse(time.RFC3339, reservation.ReservationFinishDate)
	if terr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠の取得に失敗しました", Log: terr.Error()}
	}
	finishDate, terr := time.Parse(time.RFC3339, reservation.FinishDate)
	if terr != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "予約枠の取得に失敗しました", Log: terr.Error()}
	}

	now := time.Now()
	if now.Before(reservationStartDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約受付開始前のため予約の申込・キャンセルはできません", Log: "reservation period has not started"}
	}
	if now.After(reservationFinishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "予約受付期間が終了したため予約の申込・キャンセルはできません", Log: "reservation period has ended"}
	}
	if now.After(finishDate) {
		return &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "対象の予約枠が終了したため予約の申込・キャンセルはできません", Log: "reservation slot has already finished"}
	}

	return nil
}
