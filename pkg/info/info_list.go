package info

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ResponseGetInfoList struct {
	Infos []InfoDetail `json:"infos"`
}

type InfoDetail struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Get info list
// @Accept json
// @Security Authorization
// @Router /info/ [get]
// @Param pages query int false "pages"
// @Success 200 {object} ResponseGetInfoList
func (c Context) GetInfoList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), title, updated_at FROM info_boards ORDER BY updated_at DESC LIMIT 100 OFFSET ?", pagesNum)
	if err != nil {
		e.JSON(http.StatusOK, Error{Message: "お知らせ一覧の取得に失敗しました"})
	}
	defer rows.Close()
	infos := []InfoDetail{}
	for rows.Next() {
		info := InfoDetail{}
		if err := rows.Scan(&info.ID, &info.Title, &info.UpdatedAt); err != nil {
			return e.JSON(http.StatusInternalServerError, Error{Message: "お知らせ一覧の取得に失敗しました"})
		}
		infos = append(infos, info)
	}
	return e.JSON(http.StatusOK, ResponseGetInfoList{Infos: infos})
}
