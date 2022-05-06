package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type ResponseJoin struct {
	Error string `json:"error"`
}

// Join group
// @Accept json
// @Security Authorization
// @Router /group/{id} [POST]
// @Param id path string true "group id"
// @Success 200 {object} ResponseJoin
func (c Context) Join(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGroupList{Error: err.Error()})
	}
	id := e.Param("id")
	join := true
	err = c.DB.QueryRow("SELECT `join` FROM `Group` WHERE id = UUID_TO_BIN(?)", id).Scan(&join)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseJoin{Error: err.Error()})
	}
	if !join {
		return e.JSON(http.StatusForbidden, ResponseJoin{Error: "参加権限がありません"})
	}
	_, err = c.DB.Exec("INSERT INTO GroupUser (group_id, user_id) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?))", id, userId)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return e.JSON(http.StatusForbidden, ResponseJoin{Error: "参加済みです"})
	}
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseJoin{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseJoin{Error: ""})
}