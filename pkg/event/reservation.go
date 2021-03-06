package event

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type RequestReservation struct {
	Comment string `json:"comment"`
	Url     string `json:"url"`
}

func (p RequestReservation) validate() error {
	errorMsg := []string{}
	if 255 < utf8.RuneCountInString(p.Comment) {
		errorMsg = append(errorMsg, "コメントは255文字未満である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.Url) {
		errorMsg = append(errorMsg, "URLは255文字未満である必要があります")
	}
	if len(errorMsg) != 0 {
		return fmt.Errorf(strings.Join(errorMsg, ","))
	}
	return nil
}

type ResponseReservation struct {
	Error string `json:"error"`
}

type ResponseReservationInfo struct {
	ReservationUsers []ReservationUser `json:"reservation_users"`
	Error            string            `json:"error"`
}

type ReservationUser struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Url     string `json:"url"`
}

type ResponseCancelReservation struct {
	Error string `json:"error"`
}

// Reservation event
// @Router /event/{event_id}/{id} [post]
// @Param event_id path string true "event id"
// @Param id path string true "reservation id"
// @Param RequestReservation body RequestReservation true "reservation info"
// @Security Authorization
// @Success 200 {object} ResponseReservation
func (c Context) Reservation(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseReservation{Error: err.Error()})
	}
	eventId := e.Param("event_id")
	id := e.Param("id")
	reservation := RequestReservation{}
	if err := e.Bind(&reservation); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseReservation{Error: "データの読み込みに失敗しました"})
	}
	if err := reservation.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseReservation{Error: err.Error()})
	}
	capacity := 0
	reservationStartDate := time.Time{}
	reservationFinishDate := time.Time{}
	err = c.DB.QueryRow("SELECT capacity, reservation_start_date, reservation_finish_date FROM event_reservations WHERE id = UUID_TO_BIN(?) AND event_id = UUID_TO_BIN(?)", id, eventId).Scan(&capacity, &reservationStartDate, &reservationFinishDate)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseReservation{Error: err.Error()})
	}
	nowDate := time.Now()
	if reservationStartDate.After(nowDate) || reservationFinishDate.Before(nowDate) {
		return e.JSON(http.StatusForbidden, ResponseReservation{Error: "DBの読み込みに失敗しました"})
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
	if capacity < reservatedCount+1 {
		return e.JSON(http.StatusForbidden, ResponseReservation{Error: "予約可能な枠がありません"})
	}
	_, err = c.DB.Exec("INSERT INTO event_reservation_users (reservation_id, user_id, comment, url) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?)", id, userId, reservation.Comment, reservation.Url)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return e.JSON(http.StatusBadRequest, ResponseReservation{Error: "予約済みです"})
	}
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseReservation{Error: "DBへの書き込みに失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseReservation{})
}

// Cancel reservation
// @Accept json
// @Security Authorization
// @Router /event/{event_id}/{id} [DELETE]
// @Param event_id path string true "event id"
// @Param id path string true "reservation id"
// @Success 200 {object} ResponseCancelReservation
func (c Context) CancelReservation(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseCancelReservation{Error: err.Error()})
	}
	id := e.Param("id")
	_, err = c.DB.Exec("DELETE FROM event_reservation_users WHERE reservation_id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)", id, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseCancelReservation{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseCancelReservation{})
}

// Reservation info
// @Router /event/{event_id}/{id} [get]
// @Param event_id path string true "event id"
// @Param id path string true "reservation id"
// @Security Authorization
// @Success 200 {object} ResponseReservationInfo
func (c Context) ReservationInfo(e echo.Context) error {
	id := e.Param("id")
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(users.id), user_profiles.username, comment, url FROM event_reservation_users LEFT JOIN users ON users.id = event_reservation_users.user_id LEFT JOIN user_profiles ON users.id = user_profiles.user_id WHERE reservation_id = UUID_TO_BIN(?)", id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	reservation_users := []ReservationUser{}
	for rows.Next() {
		reservation_user := ReservationUser{}
		if err := rows.Scan(&reservation_user.Id, &reservation_user.Name, &reservation_user.Comment, &reservation_user.Url); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
		}
		reservation_users = append(reservation_users, reservation_user)
	}
	return e.JSON(http.StatusOK, ResponseReservationInfo{ReservationUsers: reservation_users})
}
