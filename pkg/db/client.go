package db

import (
	"context"
	"database/sql"
	"embed"

	"github.com/future-architect/go-twowaysql"
)

type client struct {
	tw    *twowaysql.Twowaysql
	query *embed.FS
}

func (t *client) Select(dest interface{}, queryPath string, params interface{}) error {
	query, err := t.query.ReadFile(queryPath)
	if err != nil {
		return err
	}
	return t.tw.Select(context.Background(), dest, string(query), params)
}

func (t *client) Exec(queryPath string, params interface{}) (sql.Result, error) {
	query, err := t.query.ReadFile(queryPath)
	if err != nil {
		return nil, err
	}
	return t.tw.Exec(context.Background(), string(query), params)
}
