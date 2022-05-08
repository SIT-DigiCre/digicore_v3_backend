package group

import "database/sql"

type Context struct {
	DB *sql.DB
}

type Group struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Joinable    bool   `json:"joinable"`
	Joined      bool   `json:"joined"`
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}
