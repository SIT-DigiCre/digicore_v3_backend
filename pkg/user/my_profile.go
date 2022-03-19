package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Profile struct {
	UserId                string    `json:"id"`
	Username              string    `json:"username"`
	SchoolGrade           int       `json:"school_grade"`
	IconURL               string    `json:"icon_url"`
	DiscordUserId         string    `json:"discord_userid"`
	ActiveLimit           time.Time `json:"active_limit"`
	ShortSelfIntroduction string    `json:"short_self_introduction"`
}

func (p Profile) validate() error {
	return nil
}

type ResponseGetMyProfile struct {
	Profile Profile `json:"profile"`
	Error   string  `json:"error"`
}

type ResponseSetMyProfile struct {
	Error string `json:"error"`
}

// Get my prodile
// @Accept json
// @Router /user/my [get]
// @Success 200 {object} ResponseGetMyProfile
func (c Context) GetMyProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: ""})
	}
	profile := Profile{UserId: userId}
	err = c.DB.QueryRow("SELECT username, school_grade, icon_url, discord_userid, short_self_introduction FROM UserProfile WHERE user_id = UUID_TO_BIN(?)", userId).
		Scan(&profile.Username, &profile.SchoolGrade, &profile.IconURL, &profile.DiscordUserId, &profile.ShortSelfIntroduction)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseGetMyProfile{Profile: profile})
}

// Set my prodile
// @Accept json
// @Param Profile body Profile true "my profile"
// @Router /user/my [post]
// @Success 200 {object} ResponseSetMyProfile
func (c Context) SetMyProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{})
	}
	fmt.Println(userId)
	profile := Profile{}
	if err := e.Bind(&profile); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	if err := profile.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	_, err = c.DB.Exec(`INSERT INTO UserProfile (user_id, username, school_grade, icon_url, discord_userid, short_self_introduction) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)
				ON DUPLICATE KEY UPDATE username = ?, school_grade = ?, icon_url = ?, discord_userid = ?, short_self_introduction = ?`,
		userId, profile.Username, profile.SchoolGrade, profile.IconURL, profile.DiscordUserId, profile.ShortSelfIntroduction,
		profile.Username, profile.SchoolGrade, profile.IconURL, profile.DiscordUserId, profile.ShortSelfIntroduction)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseSetMyPrivateProfile{})
}
