package group

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResponseDetail struct {
	Group Group  `json:"group"`
	Error string `json:"error"`
}

// Get group detail
// @Accept json
// @Security Authorization
// @Router /group/{id} [GET]
// @Param id path string true "group id"
// @Success 200 {object} ResponseDetail
func (c Context) Detail(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseDetail{Error: err.Error()})
	}
	id := e.Param("id")
	group := Group{}
	err = c.DB.QueryRow("SELECT BIN_TO_UUID(Group.id), name, description, `join`, (CASE WHEN GroupUser.user_id IS NOT NULL THEN true ELSE false END) AS joined  FROM `Group` LEFT JOIN GroupUser ON Group.id = GroupUser.group_id AND GroupUser.user_id = UUID_TO_BIN(?) WHERE Group.id = UUID_TO_BIN(?)", userId, id).
		Scan(&group.Id, &group.Name, &group.Description, &group.Join, &group.Joined)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusBadRequest, ResponseDetail{Error: "グループが存在しません"})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseDetail{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseDetail{Group: group})
}
