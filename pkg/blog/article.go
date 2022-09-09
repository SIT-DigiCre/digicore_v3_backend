package blog

import (
	"net/http"
	"strconv"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Article struct {
	Id	string	`json:"id"`
	UserId	string	`json:"user_id"`
	Title	string	`json:"title"`
	Body	string	`json:"body"`
}

type RequestCreateArticle struct {
	Title 	string 	`json:"title"`
	Body 	string	`json:"body"`
}

type ResponseCreateArticle struct {
	ID		string	`json:"id"`
	Error	string	`json:"error"`
}

type ArticleItem struct {
	Id      string  `json:"id"`
	UserId  string  `json:"user_id"`
	Title   string  `json:"title"`
	Body    string  `json:"body"`
}

type ResponseArticleList struct {
	Error       string          `json:"error"`
	Articles    []ArticleItem   `json:"articles"`
}

func (c Context) CreateArticle(e echo.Context) error {
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseCreateArticle{Error: err.Error()})
	}
	var postCreate RequestCreateArticle
	if err := e.Bind(&postCreate); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseCreateArticle{Error: "データの読み込みに失敗しました"})
	}
	id := uuid.New().String()
	_, err = c.DB.Exec("INSERT INTO blog_posts (id, user_id, title, body) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?)", id, userId, postCreate.Title, postCreate.Body)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseCreateArticle{Error: "データベースの書き込み中にエラーが発生しました:\n" + err.Error()})	//	TODO: err.Error()は検証用なので消しておく
	}
	return e.JSON(http.StatusOK, ResponseCreateArticle{ID: id})
}

func (c Context) GetArticleList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	seed := e.QueryParam("seed")
	seedNum, _ := strconv.Atoi(seed)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), BIN_TO_UUID(user_id), title, body FROM `blog_posts` ORDER BY rand(?) LIMIT 100 OFFSET ?", seedNum, pagesNum)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	articles := []ArticleItem{}
	for rows.Next() {
		article := ArticleItem{}
		if err := rows.Scan(&article.Id, &article.UserId, &article.Title, &article.Body); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
		}
		articles = append(articles, article)
	}
	return e.JSON(http.StatusOK, ResponseArticleList{Articles: articles})
}