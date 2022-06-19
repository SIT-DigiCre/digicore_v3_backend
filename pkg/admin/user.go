package admin

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	Id                    string    `json:"id"`
	StudentNumber         string    `json:"student_number"`
	Username              string    `json:"username"`
	SchoolGrade           int       `json:"school_grade"`
	IconURL               string    `json:"icon_url"`
	DiscordUserId         string    `json:"discord_userid"`
	ActiveLimit           time.Time `json:"active_limit"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	FirstNameKana         string    `json:"first_name_kana"`
	LastNameKana          string    `json:"last_name_kana"`
	PhoneNumber           string    `json:"phone_number"`
	Address               string    `json:"address"`
	ParentName            string    `json:"parent_name"`
	ParentCellphoneNumber string    `json:"parent_cellphone_number"`
	ParentHomephoneNumber string    `json:"parent_homephone_number"`
	ParentAddress         string    `json:"parent_address"`
}

type ResponseGetUsersList struct {
	Users []User `json:"users"`
	Error string `json:"error"`
}

// Get all users
// @Router /admin/users [get]
// @Security Authorization
// @Success 200 {object} ResponseGetUsersList
func (c Context) GetUsersList(e echo.Context) error {
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(users.id), users.student_number, username, school_grade, icon_url, discord_userid, active_limit, first_name, last_name, first_name_kana, last_name_kana, phone_number, address, parent_name, parent_cellphone_number, parent_homephone_number, parent_address FROM users LEFT JOIN user_profiles ON user_profiles.user_id = users.id LEFT JOIN user_private_profiles ON user_private_profiles.user_id = users.id")
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetUsersList{Error: "DBの読み込みに失敗しました"})
	}
	users := []User{}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Id, &user.Username, &user.SchoolGrade, &user.IconURL, &user.DiscordUserId, &user.ActiveLimit, &user.FirstName, &user.LastName, &user.FirstNameKana, &user.LastNameKana, &user.PhoneNumber, &user.Address, &user.ParentName, &user.ParentCellphoneNumber, &user.ParentHomephoneNumber, &user.ParentAddress); err != nil {
			return e.JSON(http.StatusBadRequest, ResponseGetUsersList{Error: "DBの読み込みに失敗しました"})
		}
		users = append(users, user)
	}
	return e.JSON(http.StatusOK, ResponseGetUsersList{Users: users})
}
