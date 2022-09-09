package blog

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleDetail struct {
	Id      string  `json:"id"`
	UserId  string  `json:"user_id"`
	Title   string  `json:"title"`
	Body    string  `json:"body"`
}

type ResponseList struct {
	Error       string          `json:"error"`
	Articles    []ArticleDetail `json:"articles"`
}

func (c Context) GetList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	seed := e.QueryParam("seed")
	seedNum, _ := strconv.Atoi(seed)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), BIN_TO_UUID(user_id), title, body FROM `blog_posts` ORDER BY rand(?) LIMIT 100 OFFSET ?", seedNum, pagesNum)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	articles := []ArticleDetail{}
	for rows.Next() {
		article := ArticleDetail{}
		if err := rows.Scan(&article.Id, &article.UserId, &article.Title, &article.Body); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseList{Error: "DBの読み込みに失敗しました"})
		}
		articles = append(articles, article)
	}
	return e.JSON(http.StatusOK, ResponseList{Articles: articles})
}