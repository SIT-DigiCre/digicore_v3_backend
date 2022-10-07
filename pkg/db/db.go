package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
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
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	sdb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	if sdb.Ping() != nil {
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
		return transactionClient{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBでエラーが発生しました", Log: err.Error()}
	}
	return transactionClient{tx: txClient, query: &query, context: context}, nil
}

type Client interface {
	Select(dest interface{}, queryPath string, params interface{}) error
}

type TransactionClient interface {
	Select(dest interface{}, queryPath string, params interface{}) error
	Exec(queryPath string, params interface{}, generateID bool) (sql.Result, error)
	GetID() (string, error)
	DuplicateUpdate(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, error)
}
