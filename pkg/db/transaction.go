package db

import (
	"context"
	"database/sql"
	"embed"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/future-architect/go-twowaysql"
	"github.com/go-sql-driver/mysql"
)

type transactionClient struct {
	tx      *twowaysql.TwowaysqlTx
	query   *embed.FS
	context context.Context
}

func (t *transactionClient) Select(dest interface{}, queryPath string, params interface{}) error {
	query, err := t.query.ReadFile(queryPath)
	if err != nil {
		return err
	}
	return t.tx.Select(t.context, dest, string(query), params)
}

func (t *transactionClient) Exec(queryPath string, params interface{}, generateID bool) (sql.Result, error) {
	if generateID {
		_, err := t.Exec("sql/transaction/generate_id.sql", nil, false)
		return nil, err
	}
	query, err := t.query.ReadFile(queryPath)
	if err != nil {
		return nil, err
	}
	return t.tx.Exec(t.context, string(query), params)
}

func (t *transactionClient) Commit() *response.Error {
	err := t.tx.Commit()
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}

func (t *transactionClient) Rollback() error {
	return t.tx.Rollback()
}

func (t *transactionClient) GetID() (string, error) {
	id := struct {
		ID string `db:"id"`
	}{}
	err := t.Select(&id, "sql/transaction/get_id.sql", nil)
	if err != nil {
		return "", err
	}
	return id.ID, nil
}

func (t *transactionClient) DuplicateUpdate(insertQueryPath string, updateQueryPath string, params interface{}) (sql.Result, error) {
	res, err := t.Exec(insertQueryPath, params, false)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			res, err := t.Exec(updateQueryPath, params, false)
			if err != nil {
				return nil, err
			}
			return res, nil
		}
		return nil, err
	}
	return res, nil
}
