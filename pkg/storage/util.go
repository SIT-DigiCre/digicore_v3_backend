package storage

import (
	"errors"
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getS3Config() aws.Config {
	return aws.Config{
		Credentials:      credentials.NewStaticCredentials(env.WasabiAccessKey, env.WasabiSecretKey, ""),
		Endpoint:         aws.String(env.WasabiEndpoint),
		Region:           aws.String(env.WasabiRegion),
		S3ForcePathStyle: aws.Bool(true),
	}
}

func getSession() (*session.Session, error) {
	s3Config := getS3Config()
	goSession, err := session.NewSessionWithOptions(session.Options{Config: s3Config})
	if err != nil {
		return nil, errors.New("セッションの取得エラーです")
	}
	return goSession, nil
}

func getS3Client() (*s3.S3, error) {
	goSession, err := getSession()
	if err != nil {
		return nil, err
	}
	return s3.New(goSession), nil
}

func getBucketName(isPublic bool) string {
	if isPublic {
		return env.WasabiPublicBucket
	} else {
		return env.WasabiPrivateBucket
	}
}

func getFileNameFromIDandExt(fileID string, extension string) string {
	if len(extension) == 0 {
		return fileID
	}
	return fmt.Sprintf("%s.%s", fileID, extension)
}
