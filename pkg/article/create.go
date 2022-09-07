package article

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RequestCreateArticle struct {
	Title 	string 	`json:"title"`
	Body 	string	`json:"body"`
}

type ResponseCreateArticle struct {
	ID		string	`json:"id"`
	Error	string	`json:"error"`
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
