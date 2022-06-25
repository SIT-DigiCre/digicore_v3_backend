package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResponseEventsList struct {
	Events []Event `json:"events"`
	Error  string  `json:"error"`
}

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Reservable  bool   `json:"reservable"`
	Reservated  bool   `json:"reservated"`
}

// Get events list
// @Router /event [get]
// @Security Authorization
// @Success 200 {object} ResponseEventsList
func (c Context) GetEventsList(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseEventsList{Error: err.Error()})
	}
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(events.id), events.name, events.Description, IF(max(reservation_finish_date) <= NOW(),true,false) as reservable, IF(count(user_id), true, false) as reservated FROM events LEFT JOIN event_reservations ON events.id = event_reservations.event_id LEFT JOIN event_reservation_users ON event_reservations.id = event_reservation_users.reservation_id AND user_id = UUID_TO_BIN(?) GROUP BY events.id", userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	events := []Event{}
	for rows.Next() {
		event := Event{}
		if err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Reservable, &event.Reservated); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
		}
		events = append(events, event)
	}
	return e.JSON(http.StatusOK, ResponseEventsList{Events: events})
}
