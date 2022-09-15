package blog

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Article struct {
	Id          string      `json:"id"`
	UserId      string      `json:"user_id"`
	Title       string      `json:"title"`
	Body	    string      `json:"body"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	IsPublic    bool        `json:"is_public"`
}

type RequestCreateArticle struct {
	Title 	    string 	`json:"title"`
	Body 	    string	`json:"body"`
	IsPublic    bool    `json:"is_public"`
}

type ResponseCreateArticle struct {
	ID		string	`json:"id"`
	Error	string	`json:"error"`
}

type ArticleItem struct {
	Id			string		`json:"id"`
	UserId		string		`json:"user_id"`
	Title		string		`json:"title"`
	IsPublic	bool		`json:"is_public"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type ResponseArticleList struct {
	Error       string          `json:"error"`
	Articles    []ArticleItem   `json:"articles"`
}

type ResponseGetArticle struct {
	Error	string	`json:"error"`
	Article	Article	`json:"article"`
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
	published_at := time.Now()
	if postCreate.IsPublic == true {
		_, err = c.DB.Exec("INSERT INTO blog_posts (id, user_id, title, body, is_public, published_at) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, true, ?)", id, userId, postCreate.Title, postCreate.Body, published_at)
	} else {
		_, err = c.DB.Exec("INSERT INTO blog_posts (id, user_id, title, body, is_public, published_at) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, false, ?)", id, userId, postCreate.Title, postCreate.Body, nil)
	}
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseCreateArticle{Error: "データベースの書き込み中にエラーが発生しました:\n" + err.Error()})	//	TODO: err.Error()は検証用なので消しておく
	}
	return e.JSON(http.StatusOK, ResponseCreateArticle{ID: id})
}

func (c Context) GetArticleList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), BIN_TO_UUID(user_id), title, is_public, created_at, updated_at FROM `blog_posts` WHERE is_public = true ORDER BY published_at LIMIT 100 OFFSET ?", pagesNum)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	articles := []ArticleItem{}
	for rows.Next() {
		article := ArticleItem{}
		if err := rows.Scan(&article.Id, &article.UserId, &article.Title, &article.IsPublic, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
		}
		articles = append(articles, article)
	}
	return e.JSON(http.StatusOK, ResponseArticleList{Articles: articles})
}

func (c Context) GetArticle(e echo.Context) error {
	id := e.Param("id")
	article := Article{Id: id}
	err := c.DB.QueryRow("SELECT BIN_TO_UUID(user_id), title, body, created_at, updated_at, is_public FROM blog_posts WHERE id = UUID_TO_BIN(?)", id).
		Scan(&article.UserId, &article.Title, &article.Body, &article.CreatedAt, &article.UpdatedAt, &article.IsPublic)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetArticle{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetArticle{Error: "取得に失敗しました"})
	} 
	u, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseCreateArticle{Error: err.Error()})
	} else if !article.IsPublic && u != article.UserId {
		return e.JSON(http.StatusForbidden, ResponseGetArticle{Error: "アクセスが許可されていません"})
	}
	return e.JSON(http.StatusOK, ResponseGetArticle{Article: article})
}
