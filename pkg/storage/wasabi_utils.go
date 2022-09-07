package storage

import (
	"fmt"
	"bytes"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
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

func createPutObjectInput(data []byte, key string, isPublic bool) *s3.PutObjectInput {
	bucketName := getBucketName(isPublic)
	return &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}
}

func getFileURL(key string, isPublic bool) (string, error) {
	bucketName := getBucketName(isPublic)
	if isPublic {
		return fmt.Sprintf("https://%s/%s/%s", env.WasabiDirectURLDomain, bucketName, key), nil
	}
	s3Client, err := getS3Client()
	if err != nil {
		return "", err
	}
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	url, err := req.Presign(time.Hour)
	if err != nil {
		return "", errors.New("Pre-signed URL発行エラーです")
	}
	return url, nil
}
