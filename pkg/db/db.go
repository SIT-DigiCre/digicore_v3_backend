package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/future-architect/go-twowaysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var tw *twowaysql.Twowaysql

//go:embed sql
var query embed.FS

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
		env.Get("DB_USER"), env.Get("DB_PASSWORD"), env.Get("DB_HOST"), env.Get("DB_PORT"), env.Get("DB_DATABASE"))
	sdb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	err = sdb.Ping()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	tw = twowaysql.New(sdb)
}

func Open() client {
	return client{tw: tw, query: &query}
}

func OpenTransaction() (transactionClient, *response.Error) {
	context := context.Background()
	txClient, err := tw.Begin(context)
	if err != nil {
		return transactionClient{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return transactionClient{tx: txClient, query: &query, context: context}, nil
}

type Client interface {
	Select(dest interface{}, queryPath string, params interface{}) error
}

type TransactionClient interface {
	Select(dest interface{}, queryPath string, params interface{}) error
	Exec(queryPath string, params interface{}, generateId bool) (sql.Result, error)
	GetId() (string, error)
	DuplicateUpdate(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, bool, error)
}
