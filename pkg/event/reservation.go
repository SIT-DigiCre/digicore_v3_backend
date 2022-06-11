package event

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type ResponseReservation struct {
	Events []Event `json:"event"`
	Error  string  `json:"error"`
}

// Reservation event
// @Router /event/{event_id}/{id} [post]
// @Param event_id path string true "event id"
// @Param id path string true "reservation id"
// @Security Authorization
// @Success 200 {object} ResponseReservation
func (c Context) Reservation(e echo.Context) error {

	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseReservation{Error: err.Error()})
	}
	eventId := e.Param("event_id")
	id := e.Param("id")
	reservationLimit := 0
	startDate := time.Time{}
	finishDate := time.Time{}
	err = c.DB.QueryRow("SELECT reservation_limit, start_date, finish_date FROM event_reservations WHERE id = UUID_TO_BIN(?) AND event_id = UUID_TO_BIN(?)", id, eventId).Scan(&reservationLimit, &startDate, &finishDate)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseReservation{Error: "DBの読み込みに失敗しました"})
	}
	nowDate := time.Now()
	if startDate.Before(nowDate) || finishDate.After(nowDate) {
		return e.JSON(http.StatusForbidden, ResponseReservation{Error: "予約期間外です"})
	}
	myReservatedCount := 0
	err = c.DB.QueryRow("SELECT count(*) FROM events LEFT JOIN event_reservations ON events.id = event_reservations.event_id LEFT JOIN event_reservation_users ON event_reservations.id = event_reservation_users.reservation_id WHERE events.id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)", eventId, userId).Scan(&myReservatedCount)
	if err != nil || myReservatedCount != 0 {
		return e.JSON(http.StatusForbidden, ResponseReservation{Error: "予約済みです"})
	}
	reservatedCount := 0
	err = c.DB.QueryRow("SELECT count(*) FROM event_reservation_users WHERE id = UUID_TO_BIN(?)", id).Scan(&reservatedCount)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseReservation{Error: "DBの読み込みに失敗しました"})
	}
	if reservatedCount < reservationLimit {
		return e.JSON(http.StatusForbidden, ResponseReservation{Error: "予約可能な枠がありません"})
	}
	_, err = c.DB.Exec("INSERT INTO event_reservation_users (reservation_id, user_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?))", id, userId)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return e.JSON(http.StatusBadRequest, ResponseReservation{Error: "予約済みです"})
	}
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseReservation{Error: "DBへの書き込みに失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseReservation{})
}
