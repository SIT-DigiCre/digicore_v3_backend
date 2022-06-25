package event

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResponseEventDetail struct {
	Id                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Reservated        bool               `json:"reservated"`
	ReservationFrames []ReservationFrame `json:"reservation_frames"`
	Error             string             `json:"error"`
}

type ReservationFrame struct {
	Id                    string    `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	StartDate             time.Time `json:"start_date"`
	FinishDate            time.Time `json:"finish_date"`
	ReservationStartDate  time.Time `json:"reservation_start_date"`
	ReservationFinishDate time.Time `json:"reservation_finish_date"`
	Capacity              int       `json:"capacity"`
	FreeCapacity          int       `json:"free_capacity"`
	Reservable            bool      `json:"reservable"`
	Reservated            bool      `json:"reservated"`
}

// Get event detail
// @Router /event/{id} [GET]
// @Param id path string true "group id"
// @Security Authorization
// @Success 200 {object} ResponseEventDetail
func (c Context) GetEventDetail(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseEventsList{Error: err.Error()})
	}
	id := e.Param("id")
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(event_reservations.id), event_reservations.name, event_reservations.Description,start_date, finish_date, reservation_start_date, reservation_finish_date, event_reservations.capacity, event_reservations.capacity - count(event_reservation_users.id), IF(SUM(event_reservation_users.user_id = UUID_TO_BIN(?)) ,true ,false) AS reservated  FROM event_reservations LEFT JOIN event_reservation_users ON event_reservations.event_id = UUID_TO_BIN(?) AND event_reservations.id = event_reservation_users.reservation_id GROUP BY event_reservations.id", userId, id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: err.Error()})
	}
	defer rows.Close()
	reservation_frames := []ReservationFrame{}
	for rows.Next() {
		reservation_frame := ReservationFrame{}
		if err := rows.Scan(&reservation_frame.Id, &reservation_frame.Name, &reservation_frame.Description, &reservation_frame.StartDate, &reservation_frame.FinishDate, &reservation_frame.ReservationStartDate, &reservation_frame.ReservationFinishDate, &reservation_frame.Capacity, &reservation_frame.FreeCapacity, &reservation_frame.Reservated); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
		}
		now_date := time.Now()
		if reservation_frame.FreeCapacity == 0 || now_date.Before(reservation_frame.ReservationStartDate) || now_date.After(reservation_frame.ReservationFinishDate) {
			reservation_frame.Reservable = false
		} else {
			reservation_frame.Reservable = true
		}
		reservation_frames = append(reservation_frames, reservation_frame)
	}
	response := ResponseEventDetail{ReservationFrames: reservation_frames}
	err = c.DB.QueryRow("SELECT BIN_TO_UUID(events.id), events.name, events.description, (CASE WHEN user_id IS NOT NULL THEN true ELSE false END) AS reservated FROM events LEFT JOIN event_reservations ON events.id = event_reservations.event_id LEFT JOIN event_reservation_users ON event_reservations.id = event_reservation_users.reservation_id AND event_reservation_users.user_id = UUID_TO_BIN(?) WHERE events.id = UUID_TO_BIN(?)", userId, id).Scan(&response.Id, &response.Name, &response.Description, &response.Reservated)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
	}
	return e.JSON(http.StatusOK, response)
}
