package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/future-architect/go-twowaysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Client interface {
	Select(dest interface{}, query string, params interface{}) error
	Exec(query string, params interface{}) (sql.Result, error)
}

type DBClient struct {
	Client Client
	Query  embed.FS
}

var DB DBClient

//go:embed sql
var query embed.FS

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	sdb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	if sdb.Ping() != nil {
		logrus.Fatal(err.Error())
	}
	DB.Client = &client{tw: twowaysql.New(sdb)}
	DB.Query = query
}

type client struct {
	tw *twowaysql.Twowaysql
}

func (c *client) Select(dest interface{}, query string, params interface{}) error {
	return c.tw.Select(context.Background(), dest, query, params)
}

func (c *client) Exec(query string, params interface{}) (sql.Result, error) {
	return c.tw.Exec(context.Background(), query, params)
}
