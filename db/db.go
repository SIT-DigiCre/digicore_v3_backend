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

var (
	Client DBClient
	//go:embed sql/*
	Query embed.FS
)

type DBClient interface {
	Select(ctx context.Context, dest interface{}, query string, params interface{}) error
	Exec(ctx context.Context, query string, params interface{}) (sql.Result, error)
	Close() error
}

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	sdb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logrus.Error(err.Error())
	}
	Client = &db{tw: twowaysql.New(sdb)}
}

type db struct {
	tw *twowaysql.Twowaysql
}

func (d *db) Select(ctx context.Context, dest interface{}, query string, params interface{}) error {
	return d.tw.Select(ctx, dest, query, params)
}

func (d *db) Exec(ctx context.Context, query string, params interface{}) (sql.Result, error) {
	return d.tw.Exec(ctx, query, params)
}

func (d *db) Close() error {
	return d.tw.Close()
}
