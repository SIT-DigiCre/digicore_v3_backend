package event

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseEventDetail struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Detail []Detail `json:"detail"`
	Error  string   `json:"error"`
}

type Detail struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Full        bool   `json:"full"`
}

// Get event detail
// @Router /event/{id} [GET]
// @Param id path string true "group id"
// @Security Authorization
// @Success 200 {object} ResponseEventDetail
func (c Context) GetEventDetail(e echo.Context) error {
	id := e.Param("id")
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(event_reservations.id), event_reservations.name, event_reservations.Description, IF(event_reservations.reservation_limit <= count(event_reservation_users.id) ,true ,false) AS full FROM event_reservations LEFT JOIN event_reservation_users ON event_reservations.event_id = UUID_TO_BIN(?) AND event_reservations.id = event_reservation_users.reservation_id GROUP BY event_reservations.id", id)
	defer rows.Close()
	details := []Detail{}
	for rows.Next() {
		detail := Detail{}
		if err := rows.Scan(&detail.Id, &detail.Name, &detail.Description, &detail.Full); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
		}
		details = append(details, detail)
	}
	response := ResponseEventDetail{Detail: details}
	err = c.DB.QueryRow("SELECT BIN_TO_UUID(id), name, description FROM events WHERE id = UUID_TO_BIN(?)", id).Scan(&response.Id,&response.Name,&response.Description)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseEventsList{Error: "DBの読み込みに失敗しました"})
	}
	return e.JSON(http.StatusOK,response )
}
