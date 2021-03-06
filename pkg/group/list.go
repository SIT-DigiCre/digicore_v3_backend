package group

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type ResponseList struct {
	Groups []Group `json:"groups"`
	Error  string  `json:"error"`
}

// Get group list
// @Accept json
// @Security Authorization
// @Router /group [get]
// @Success 200 {object} ResponseList
func (c Context) GroupList(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseList{Error: err.Error()})
	}
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(groups.id), name, description, joinable, (CASE WHEN groups_users.user_id IS NOT NULL THEN true ELSE false END) AS joined  FROM `groups` LEFT JOIN groups_users ON groups.id = groups_users.group_id AND groups_users.user_id = UUID_TO_BIN(?)", userId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	var groups []Group
	for rows.Next() {
		group := Group{}
		if err := rows.Scan(&group.Id, &group.Name, &group.Description, &group.Joinable, &group.Joined); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseList{Error: "DBの読み込みに失敗しました"})
		}
		groups = append(groups, group)
	}
	return e.JSON(http.StatusOK, ResponseList{Groups: groups})
}
