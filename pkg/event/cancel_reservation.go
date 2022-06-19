package event

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResponseCancelReservation struct {
	Error string `json:"error"`
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
