package blog

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Article struct {
	Id          string      	`json:"id"`
	UserId      string      	`json:"user_id"`
	Title       string      	`json:"title"`
	Body	    string      	`json:"body"`
	IsPublic    bool        	`json:"is_public"`
	PublishedAt	time.Time		`json:"published_at"`
	CreatedAt   time.Time   	`json:"created_at"`
	UpdatedAt   time.Time   	`json:"updated_at"`
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
	Id			string			`json:"id"`
	UserId		string			`json:"user_id"`
	Title		string			`json:"title"`
	IsPublic	bool			`json:"is_public"`
	PublishedAt	time.Time		`json:"published_at"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
}

type ResponseArticleList struct {
	Error       string          `json:"error"`
	Articles    []ArticleItem   `json:"articles"`
}

type ResponseGetArticle struct {
	Error	string	`json:"error"`
	Article	Article	`json:"article"`
}

type RequestUpdateArticle struct {
	Title 	    string 	`json:"title"`
	Body 	    string	`json:"body"`
	IsPublic    bool    `json:"is_public"`
}

func (a RequestUpdateArticle) validate() error {
	errorMsg := []string{}
	if utf8.RuneCountInString(a.Title) < 1 {
		errorMsg = append(errorMsg, "タイトルは1文字以上である必要があります")
	}
	if utf8.RuneCountInString(a.Body) < 1 {
		errorMsg = append(errorMsg, "本文は1文字以上である必要があります")
	}
	if len(errorMsg) != 0 {
		return fmt.Errorf(strings.Join(errorMsg, ","))
	}
	return nil
}

type ResponseUpdateArticle struct {
	Error	string	`json:"error"`
}

type ResponseDeleteArticle struct {
	Error	string	`json:"error"`
}

//	Create article
//	@Accept json
//	@Param RequestCreateArticle body RequestCreateArticle true "article data"
//	@Security Authorization
//	@Router /blog/article [post]
//	@Success 200 {object} ResponseCreateArticle
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
	published_at := "2000-01-01T00:00:00+00:00"
	if postCreate.IsPublic {
		published_at = "CURRENT_TIMESTAMP"
	}
	_, err = c.DB.Exec("INSERT INTO blog_posts (id, user_id, title, body, is_public, published_at) VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?)", id, userId, postCreate.Title, postCreate.Body, postCreate.IsPublic, published_at)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseCreateArticle{Error: "データベースの書き込み中にエラーが発生しました:\n" + err.Error()})	//	TODO: err.Error()は検証用なので消しておく
	}
	return e.JSON(http.StatusOK, ResponseCreateArticle{ID: id})
}

//	Get public article list
//	@Router /blog/articles [get]
//	@Param pages query int false "pages"
//	@Security Authorization
//	@Success 200 {object} ResponseArticleList
func (c Context) GetArticleList(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), BIN_TO_UUID(user_id), title, is_public, published_at, created_at, updated_at FROM `blog_posts` WHERE is_public = true ORDER BY published_at LIMIT 100 OFFSET ?", pagesNum)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	articles := []ArticleItem{}
	for rows.Next() {
		article := ArticleItem{}
		if err := rows.Scan(&article.Id, &article.UserId, &article.Title, &article.IsPublic, &article.PublishedAt, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
		}
		articles = append(articles, article)
	}
	return e.JSON(http.StatusOK, ResponseArticleList{Articles: articles})
}

//	Get list of my articles
//	@Router /blog/articles/my [get]
//	@Param pages query int false "pages"
//	@Security Authorization
//	@Success 200 {object} ResponseArticleList
func (c Context) GetMyArticles(e echo.Context) error {
	pages := e.QueryParam("pages")
	pagesNum, _ := strconv.Atoi(pages)
	userId, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseArticleList{Error: err.Error()})
	}
	rows, err := c.DB.Query("SELECT BIN_TO_UUID(id), title, is_public, published_at, created_at, updated_at FROM `blog_posts` WHERE user_id = UUID_TO_BIN(?) ORDER BY published_at LIMIT 100 OFFSET ?", userId, pagesNum)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
	}
	defer rows.Close()
	articles := []ArticleItem{}
	for rows.Next() {
		article := ArticleItem{UserId: userId}
		if err := rows.Scan(&article.Id, &article.Title, &article.IsPublic, &article.PublishedAt, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return e.JSON(http.StatusInternalServerError, ResponseArticleList{Error: "DBの読み込みに失敗しました"})
		}
		articles = append(articles, article)
	}
	return e.JSON(http.StatusOK, ResponseArticleList{Articles: articles})
}

//	Get article
//	@Router /blog/articles/{id} [get]
//	@Param id path string true "article id"
//	@Security Authorization
//	@Success 200 {object} ResponseGetArticle
func (c Context) GetArticle(e echo.Context) error {
	id := e.Param("id")
	article := Article{Id: id}
	err := c.DB.QueryRow("SELECT BIN_TO_UUID(user_id), title, body, is_public, published_at, created_at, updated_at FROM blog_posts WHERE id = UUID_TO_BIN(?)", id).
		Scan(&article.UserId, &article.Title, &article.Body, &article.IsPublic, &article.PublishedAt, &article.CreatedAt, &article.UpdatedAt)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseGetArticle{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseGetArticle{Error: "取得に失敗しました"})
	} 
	u, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseGetArticle{Error: err.Error()})
	} else if !article.IsPublic && u != article.UserId {
		return e.JSON(http.StatusForbidden, ResponseGetArticle{Error: "アクセスが許可されていません"})
	}
	return e.JSON(http.StatusOK, ResponseGetArticle{Article: article})
}

//	Update article
//	@Accept json
//	@Router /blog/articles/{id} [put]
//	@Param id path string true "article id"
//	@Param RequestUpdateArticle body RequestUpdateArticle true "article data"
//	@Security Authorization
//	@Success 200 {object} ResponseUpdateArticle
func (c Context) UpdateArticle(e echo.Context) error {
	id := e.Param("id")
	article := RequestUpdateArticle{}
	if err := e.Bind(&article); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateArticle{Error: "データの読み込みに失敗しました"})
	}
	if err := article.validate(); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateArticle{Error: err.Error()})
	}
	original := Article{}
	err := c.DB.QueryRow("SELECT BIN_TO_UUID(user_id), published_at FROM blog_posts WHERE id = UUID_TO_BIN(?)", id).
		Scan(&original.UserId, &original.PublishedAt)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseUpdateArticle{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseUpdateArticle{Error: "取得に失敗しました"})
	}
	u, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateArticle{Error: err.Error()})
	} else if !article.IsPublic && u != original.UserId {
		return e.JSON(http.StatusForbidden, ResponseUpdateArticle{Error: "アクセスが許可されていません"})
	}
	if article.IsPublic && original.PublishedAt.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		_, err = c.DB.Exec(`UPDATE blog_posts SET title = ?, body = ?, is_public = ?, published_at = CURRENT_TIMESTAMP WHERE id = UUID_TO_BIN(?)`,
			article.Title, article.Body, article.IsPublic, id)
	} else {
		_, err = c.DB.Exec(`UPDATE blog_posts SET title = ?, body = ?, is_public = ? WHERE id = UUID_TO_BIN(?)`,
			article.Title, article.Body, article.IsPublic, id)
	}
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseUpdateArticle{Error: "更新に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdateArticle{})
}

//	Delete article
//	@Router /blog/articles/{id} [delete]
//	@Param id path string true "article id"
//	@Security Authorization
//	@Success 200 {object} ResponseDeleteArticle
func (c Context) DeleteArticle(e echo.Context) error {
	id := e.Param("id")
	article := Article{}
	err := c.DB.QueryRow("SELECT BIN_TO_UUID(user_id) FROM blog_posts WHERE id = UUID_TO_BIN(?)", id).
		Scan(&article.UserId)
	if err == sql.ErrNoRows {
		return e.JSON(http.StatusNotFound, ResponseDeleteArticle{Error: "データが登録されていません"})
	} else if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseDeleteArticle{Error: "取得に失敗しました"})
	}
	u, err := user.GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseDeleteArticle{Error: err.Error()})
	} else if !article.IsPublic && u != article.UserId {
		return e.JSON(http.StatusForbidden, ResponseDeleteArticle{Error: "アクセスが許可されていません"})
	}
	_, err = c.DB.Exec(`DELETE FROM blog_posts WHERE id = UUID_TO_BIN(?)`, id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseDeleteArticle{Error: "削除に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseDeleteArticle{})
}