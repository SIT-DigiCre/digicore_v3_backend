package user

import (
	"database/sql"
)

type Context struct {
	DB *sql.DB
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}
