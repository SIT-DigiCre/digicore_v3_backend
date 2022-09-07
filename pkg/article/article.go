package article

import (
	"database/sql"
)

type Context struct {
	DB *sql.DB
}

type Blog struct {
	Id	string	`json:"id"`
	UserId	string	`json:"user_id"`
	Title	string	`json:"title"`
	Body	string	`json:"body"`
}

func CreateContext(db *sql.DB) (Context, error){
	context := Context{DB: db}

	return context, nil
}
