package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	if os.Getenv("AWS_SAM_LOCAL") != "" {
	}
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, tweets []Tweet) (string, error) {
	spew.Dump(tweets)
	// S3 Client test
	sess := session.New(&aws.Config{})
	s3Svc := s3.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))
	uploader := s3manager.NewUploaderWithClient(s3Svc)
	contentType := "application/json"
	bucket := os.Getenv("BucketName")
	key := "hogehogeeeeetests"
	jsonBytes, err := json.Marshal(tweets)
	if err != nil {
		return "", err
	}
	upParams := &s3manager.UploadInput{
		Bucket:      &bucket,
		Key:         &key,
		ACL:         aws.String("public-read"),
		Body:        bytes.NewReader(jsonBytes),
		ContentType: &contentType,
	}

	result, err := uploader.Upload(upParams)
	if err != nil {
		log.Println(result)
		log.Println(err)
		return "", err
	}
	return "s3upload.go", nil
}
