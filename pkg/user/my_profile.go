package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

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

type RequestUpdateMyProfile struct {
	Username              string `json:"username"`
	SchoolGrade           int    `json:"school_grade"`
	IconURL               string `json:"icon_url"`
	ShortSelfIntroduction string `json:"short_self_introduction"`
}

func (p RequestUpdateMyProfile) validate() error {
	errorMsg := []string{}
	if 255 < utf8.RuneCountInString(p.Username) {
		errorMsg = append(errorMsg, "ユーザー名は255文字未満である必要があります")
	}
	if p.SchoolGrade <= 0 || 10 <= p.SchoolGrade {
		errorMsg = append(errorMsg, "学年は1以上9以下である必要があります")
	}
	if 255 < utf8.RuneCountInString(p.ShortSelfIntroduction) {
		errorMsg = append(errorMsg, "自己紹介は255文字未満である必要があります")
	}
	if len(errorMsg) != 0 {
		return fmt.Errorf(strings.Join(errorMsg, ","))
	}
	return nil
}

type ResponseGetMyProfile struct {
	Profile Profile `json:"profile"`
	Error   string  `json:"error"`
}

type ResponseUpdateMyProfile struct {
	Error string `json:"error"`
}

// Get my prodile
// @Router /user/my [get]
// @Security Authorization
// @Success 200 {object} ResponseGetMyProfile
func (c Context) GetMyProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	}
	profile := Profile{UserId: userId}
	err = c.DB.QueryRow("SELECT username, school_grade, icon_url, discord_userid, active_limit, short_self_introduction FROM UserProfile WHERE user_id = ?", userId).
		Scan(&profile.Username, &profile.SchoolGrade, &profile.IconURL, &profile.DiscordUserId, &profile.ActiveLimit, &profile.ShortSelfIntroduction)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusOK, ResponseGetMyProfile{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: "取得に失敗しました"})
	}
	profile.StudentNumber, err = GetStudentNumber(c.DB, profile.UserId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseGetMyProfile{Profile: profile})
}

// Update my prodile
// @Accept json
// @Param RequestUpdateMyProfile body RequestUpdateMyProfile true "my profile"
// @Security Authorization
// @Router /user/my [put]
// @Success 200 {object} ResponseUpdateMyProfile
func (c Context) UpdateMyProfile(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetMyPrivateProfile{Error: err.Error()})
	}
	fmt.Println(userId)
	profile := RequestUpdateMyProfile{}
	if err := e.Bind(&profile); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyProfile{Error: "データの読み込みに失敗しました"})
	}
	if err := profile.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyProfile{Error: err.Error()})
	}
	_, err = c.DB.Exec(`UPDATE UserProfile SET username = ?, school_grade = ?, icon_url = ?, short_self_introduction = ? WHERE user_id = ?`,
		profile.Username, profile.SchoolGrade, profile.IconURL, profile.ShortSelfIntroduction, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateMyProfile{Error: "更新に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdateMyProfile{})
}
