package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResponseLeave struct {
	Error string `json:"error"`
}

// Get group list
// @Accept json
// @Security Authorization
// @Router /group/{id} [DELETE]
// @Param id path string true "group id"
// @Success 200 {object} ResponseLeave
func (c Context) Leave(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGroupList{Error: err.Error()})
	}
	id := e.Param("id")
	_, err = c.DB.Exec("DELETE FROM GroupUser WHERE group_id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)", id, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseJoin{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseJoin{Error: ""})
}
