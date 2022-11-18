package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func PostStorageMyfile(ctx echo.Context, dbTransactionClient db.TransactionClient, requestBody api.ReqPostStorageMyfile) (api.ResGetStorageFileId, *response.Error) {
	userId := ctx.Get("user_id").(string)
	data, rerr := base64.StdEncoding.DecodeString(requestBody.File)
	if rerr != nil {
		return api.ResGetStorageFileId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ファイルのデコードに失敗しました", Log: rerr.Error()}
	}
	fileId, err := createUserFile(dbTransactionClient, userId, requestBody, data)
	if err != nil {
		return api.ResGetStorageFileId{}, err
	}
	return GetStorageFileId(ctx, dbTransactionClient, fileId)
}

func createUserFile(dbTransactionClient db.TransactionClient, userId string, requestBody api.ReqPostStorageMyfile, data []byte) (string, *response.Error) {
	extension := getExtension(requestBody.Name)
	md5Hash := fmt.Sprintf("%x", md5.Sum(data))
	kSize := len(data) / 1024
	params := struct {
		UserId    string `twowaysql:"userId"`
		Name      string `twowaysql:"name"`
		KSize     int    `twowaysql:"kSize"`
		Md5Hash   string `twowaysql:"md5Hash"`
		Extension string `twowaysql:"extension"`
		IsPublic  bool   `twowaysql:"isPublic"`
	}{
		UserId:    userId,
		Name:      requestBody.Name,
		KSize:     kSize,
		Extension: extension,
		Md5Hash:   md5Hash,
		IsPublic:  requestBody.IsPublic,
	}
	_, rerr := dbTransactionClient.Exec("sql/storage/insert_user_storage.sql", &params, true)
	if rerr != nil {
		if mysqlErr, ok := rerr.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "既にアップロード済みのファイルです", Log: rerr.Error()}
			}
		}
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	fileId, rerr := dbTransactionClient.GetId()
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: rerr.Error()}
	}
	err := uploadBytes(data, fileId, extension, requestBody.IsPublic)
	if err != nil {
		return "", err
	}
	return fileId, nil
}

func uploadBytes(data []byte, fileId string, extension string, isPublic bool) *response.Error {
	s3Client, err := getS3Client()
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	key := getFileNameFromIDandExt(fileId, extension)
	putObjectInput := createPutObjectInput(data, key, isPublic)
	_, err = s3Client.PutObject(putObjectInput)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ファイルのアップロードに失敗しました", Log: err.Error()}
	}
	return nil
}

func createPutObjectInput(data []byte, key string, isPublic bool) *s3.PutObjectInput {
	bucketName := getBucketName(isPublic)
	return &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
}

func getExtension(fileName string) string {
	dot_ext := filepath.Ext(fileName)
	if len(dot_ext) == 0 {
		return ""
	}
	return dot_ext[1:]
}
