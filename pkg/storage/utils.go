package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
)

type Context struct {
	DB *sql.DB
}

func CreateContext(db *sql.DB) (Context, error) {
	context := Context{DB: db}

	return context, nil
}

func getExtension(fileName string) string {
	dot_ext := filepath.Ext(fileName)
	if len(dot_ext) == 0 {
		return ""
	}
	return dot_ext[1:]
}
func getFileNameFromIDandExt(fileID string, extension string) string {
	if len(extension) == 0 {
		return fileID
	}
	return fmt.Sprintf("%s.%s", fileID, extension)
}

func getFileURL(fileId string, extension string) string {
	url := fmt.Sprintf(
		"%s/%s/%s",
		env.ConohaObjectStorageServerURL,
		env.ConohaStorageContainerName,
		fileId,
	)
	if len(extension) == 0 {
		return url
	}
	return fmt.Sprintf("%s.%s", url, extension)
}

func httpRequest(method string, url string, body io.Reader, token *string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New("HTTPリクエスト作成エラー")
	}
	req.Header.Set("Accept", "application/json")
	if token != nil {
		req.Header.Set("X-Auth-Token", *token)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("HTTPリクエストエラー")
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return nil, errors.New(fmt.Sprintf("ステータスコードが%dでした", res.StatusCode))
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("レスポンスボディ取得エラー")
	}
	return resBody, nil
}
