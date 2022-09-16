package work

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Work struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Authers []Auther `json:"authers"`
	Tags    []Tag    `json:"tags"`
}

type ResponseGetWorkList struct {
	Works []Work `json:"works"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResponseGetTagList struct {
	Tags []Tag `json:"tags"`
}

// Get work list
// @Accept json
// @Security Authorization
// @Router /work/work [get]
// @Param pages query int false "pages"
// @Success 200 {object} ResponseGetWorkList
func (c Context) WorkList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), name FROM works ORDER BY updated_at DESC LIMIT 100 OFFSET ?", pagesNum)
	if err != nil {
		e.JSON(http.StatusOK, Error{Message: "作品一覧の取得に失敗しました"})
	}
	works := []Work{}
	for rows.Next() {
		work := Work{}
		if err := rows.Scan(&work.ID, &work.Name); err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "作品一覧の取得に失敗しました"})
		}
		authers := []Auther{}
		authers_rows, err := c.DB.Query("SELECT BIN_TO_UUID(work_users.user_id), username, icon_url FROM work_users LEFT JOIN user_profiles ON user_profiles.user_id = work_users.user_id WHERE work_id = UUID_TO_BIN(?)", work.ID)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "製作者の読み込みに失敗しました"})
		}
		for authers_rows.Next() {
			auther := Auther{}
			if err := authers_rows.Scan(&auther.ID, &auther.Name, &auther.IconURL); err != nil {
				return e.JSON(http.StatusInternalServerError, Error{Message: "製作者の読み込みに失敗しました"})
			}
			authers = append(authers, auther)
		}
		work.Authers = authers
		tags := []Tag{}
		tags_rows, err := c.DB.Query("SELECT BIN_TO_UUID(tag_id), name FROM work_work_tags LEFT JOIN work_tags ON work_tags.id = work_work_tags.tag_id WHERE work_id = UUID_TO_BIN(?)", work.ID)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "タグの読み込みに失敗しました"})
		}
		for tags_rows.Next() {
			tag := Tag{}
			if err := tags_rows.Scan(&tag.ID, &tag.Name); err != nil {
				return e.JSON(http.StatusInternalServerError, Error{Message: "タグの読み込みに失敗しました"})
			}
			tags = append(tags, tag)
		}
		work.Tags = tags
		works = append(works, work)
	}
	return e.JSON(http.StatusOK, ResponseGetWorkList{Works: works})
}

// Get tag list
// @Accept json
// @Security Authorization
// @Router /work/tag [get]
// @Param pages query int false "pages"
// @Success 200 {object} ResponseGetTagList
func (c Context) TagList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), name FROM work_tags ORDER BY updated_at DESC LIMIT 100 OFFSET ?", pagesNum)
	if err != nil {
		e.JSON(http.StatusOK, Error{Message: "タグ一覧の取得に失敗しました"})
	}
	tags := []Tag{}
	for rows.Next() {
		tag := Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "タグ一覧の取得に失敗しました"})
		}
		tags = append(tags, tag)
	}
	return e.JSON(http.StatusOK, ResponseGetTagList{Tags: tags})
}
