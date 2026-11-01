package storage

import (
	"bytes"
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadWorkThumbnail(data []byte, key, contentType string) error {
	client, err := getS3Client()
	if err != nil {
		return err
	}
	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(env.WasabiPublicBucket), Key: aws.String(key),
		Body: bytes.NewReader(data), ContentType: aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	return err
}

func WorkThumbnailURL(key string) string {
	return fmt.Sprintf("https://%s/%s/%s", env.WasabiDirectURLDomain, env.WasabiPublicBucket, key)
}
