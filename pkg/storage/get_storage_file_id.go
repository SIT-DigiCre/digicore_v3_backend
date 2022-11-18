package storage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func GetStorageFileId(ctx echo.Context, dbClient db.Client, fileId string) (api.ResGetStorageFileId, *response.Error) {
	res := api.ResGetStorageFileId{}
	file, err := getFileFromFileId(dbClient, fileId)
	if err != nil {
		return api.ResGetStorageFileId{}, err
	}
	rerr := copier.Copy(&res, &file)
	if rerr != nil {
		return api.ResGetStorageFileId{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: rerr.Error()}
	}
	key := getFileNameFromIDandExt(fileId, res.Extension)
	res.Url, err = getFileURL(key, res.IsPublic)
	if err != nil {
		return api.ResGetStorageFileId{}, err
	}
	return res, nil
}

func getFileURL(key string, isPublic bool) (string, *response.Error) {
	bucketName := getBucketName(isPublic)
	if isPublic {
		return fmt.Sprintf("https://%s/%s/%s", env.WasabiDirectURLDomain, bucketName, key), nil
	}
	s3Client, err := getS3Client()
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ファイルが有りません", Log: err.Error()}
	}
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	url, err := req.Presign(time.Hour)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "Pre-signed URL発行エラーです", Log: err.Error()}
	}
	return url, nil
}

func getFileFromFileId(dbClient db.Client, fileId string) (file, *response.Error) {
	params := struct {
		FileId string `twowaysql:"fileId"`
	}{
		FileId: fileId,
	}
	files := []file{}
	err := dbClient.Select(&files, "sql/storage/select_storage_from_file_id.sql", &params)
	if err != nil {
		return file{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(files) == 0 {
		return file{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "ファイルが有りません", Log: "no rows in result"}
	}
	return files[0], nil
}
