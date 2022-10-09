package info

import "database/sql"

type Context struct {
	DB *sql.DB
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}

type Error struct {
	Message string `json:"message"`
}
