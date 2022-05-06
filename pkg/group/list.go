package group

import (
	"log"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

type Group struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Join        bool   `json:"join"`
	Joined      bool   `json:"joined"`
}

type ResponseGroupList struct {
	Groups []Group `json:"groups"`
	Error  string  `json:"error"`
}

// Get group list
// @Accept json
// @Security Authorization
// @Router /group [get]
// @Success 200 {object} ResponseGroupList
func (c Context) GroupList(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGroupList{Error: err.Error()})
	}
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(Group.id), name, description, `join`, (CASE WHEN GroupUser.user_id IS NOT NULL THEN true ELSE false END) AS joined  FROM `Group` LEFT JOIN GroupUser ON Group.id = GroupUser.group_id AND GroupUser.user_id = UUID_TO_BIN(?)", userId)
	var groups []Group
	for rows.Next() {
		group := Group{}
		if err := rows.Scan(&group.Id, &group.Name, &group.Description, &group.Join, &group.Joined); err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	return e.JSON(http.StatusOK, ResponseGroupList{Groups: groups})
}