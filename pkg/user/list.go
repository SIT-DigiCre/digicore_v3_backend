package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProfileDetail struct {
	UserId                string `json:"id"`
	Username              string `json:"username"`
	IconURL               string `json:"icon_url"`
	ShortSelfIntroduction string `json:"short_self_introduction"`
}

type ResponseList struct {
	Profiles []ProfileDetail `json:"profiles"`
	Error    string          `json:"error"`
}

// Get user list
// @Accept json
// @Security Authorization
// @Router /user [get]
// @Param pages query int false "pages"
// @Param seed query int false "seed"
// @Success 200 {object} ResponseList
func (c Context) GetList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	seed := e.QueryParam("seed")
	seedNum, _ := strconv.Atoi(seed)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(user_id), username, icon_url, short_self_introduction FROM `user_profiles` ORDER BY rand(?) LIMIT 100 OFFSET ?", seedNum, pagesNum)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	profiles := []ProfileDetail{}
	for rows.Next() {
		profile := ProfileDetail{}
		if err := rows.Scan(&profile.UserId, &profile.Username, &profile.IconURL, &profile.ShortSelfIntroduction); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseList{Error: "DBの読み込みに失敗しました"})
		}
		profiles = append(profiles, profile)
	}
	return e.JSON(http.StatusOK, ResponseList{Profiles: profiles})
}
