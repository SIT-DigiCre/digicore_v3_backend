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
	StudentNumber         string    `json:"student_number"`
	Username              string    `json:"username"`
	SchoolGrade           int       `json:"school_grade"`
	IconURL               string    `json:"icon_url"`
	DiscordUserId         string    `json:"discord_userid"`
	ActiveLimit           time.Time `json:"active_limit"`
	ShortSelfIntroduction string    `json:"short_self_introduction"`
}

type UpdateableProfile struct {
	Username              string `json:"username"`
	SchoolGrade           int    `json:"school_grade"`
	IconURL               string `json:"icon_url"`
	ShortSelfIntroduction string `json:"short_self_introduction"`
}

func (p UpdateableProfile) validate() error {
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
	err = c.DB.QueryRow("SELECT username, school_grade, icon_url, discord_userid, active_limit, short_self_introduction FROM UserProfile WHERE user_id = UUID_TO_BIN(?)", userId).
		Scan(&profile.Username, &profile.SchoolGrade, &profile.IconURL, &profile.DiscordUserId, &profile.ActiveLimit, &profile.ShortSelfIntroduction)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	}
	profile.StudentNumber = GetStudentNumber(c, profile.UserId)
	return e.JSON(http.StatusOK, ResponseGetMyProfile{Profile: profile})
}

// Update my prodile
// @Accept json
// @Param Profile body UpdateableProfile true "my profile"
// @Router /user/my [put]
// @Success 200 {object} ResponseSetMyProfile
func (c Context) UpdateMyProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{})
	}
	fmt.Println(userId)
	profile := UpdateableProfile{}
	if err := e.Bind(&profile); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	if err := profile.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	_, err = c.DB.Exec(`UPDATE UserProfile SET username = ?, school_grade = ?, icon_url = ?, short_self_introduction = ? WHERE user_id = UUID_TO_BIN(?)`,
		profile.Username, profile.SchoolGrade, profile.IconURL, profile.ShortSelfIntroduction, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseSetMyPrivateProfile{})
}
