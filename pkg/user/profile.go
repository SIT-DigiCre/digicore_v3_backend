package user

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseGetProfile struct {
	Profile Profile `json:"profile"`
	Error   string  `json:"error"`
}

type SelfIntroduction struct {
	SelfIntroduction string `json:"self_introduction"`
}

type ResponseGetSelfIntroduction struct {
	SelfIntroduction SelfIntroduction `json:"self_introduction"`
	Error            string           `json:"error"`
}

// Get profile
// @Router /user/{id} [get]
// @Param id path string true "user id"
// @Security Authorization
// @Success 200 {object} ResponseGetProfile
func (c Context) GetProfile(e echo.Context) error {
	id := e.Param("id")
	profile := Profile{UserId: id}
	err := c.DB.QueryRow("SELECT username, school_grade, icon_url, discord_userid, active_limit, short_self_introduction FROM user_profiles WHERE user_id = UUID_TO_BIN(?)", id).
		Scan(&profile.Username, &profile.SchoolGrade, &profile.IconURL, &profile.DiscordUserId, &profile.ActiveLimit, &profile.ShortSelfIntroduction)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetProfile{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetProfile{Error: "取得に失敗しました"})
	}
	profile.StudentNumber, err = GetStudentNumber(c.DB, profile.UserId)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseGetProfile{Profile: profile})
}

// Get self introduction
// @Router /user/{id}/intro [get]
// @Param id path string true "user id"
// @Security Authorization
// @Success 200 {object} ResponseGetSelfIntroduction
func (c Context) GetSelfIntroduction(e echo.Context) error {
	id := e.Param("id")
	self_introduction := SelfIntroduction{}
	err := c.DB.QueryRow("SELECT self_introduction FROM user_profiles WHERE user_id = UUID_TO_BIN(?)", id).
		Scan(&self_introduction.SelfIntroduction)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetSelfIntroduction{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetSelfIntroduction{Error: "取得に失敗しました"})
	}
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetSelfIntroduction{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseGetSelfIntroduction{SelfIntroduction: self_introduction})
}
