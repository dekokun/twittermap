package main

import (
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	if os.Getenv("AWS_SAM_LOCAL") != "" {
	}
	lambda.Start(handleRequest)
}

func handleRequest(request events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	// S3 Client test
	sess := session.New(&aws.Config{})
	s3Svc := s3.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))
	uploader := s3manager.NewUploaderWithClient(s3Svc)
	bucket := "twittermap.dekokun.info"
	key := "hogehogeeeeetests"
	upParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		ACL:    aws.String("public-read"),
		Body:   strings.NewReader("fugafuga"),
	}

	result, err := uploader.Upload(upParams)
	log.Println(result)
	log.Println(err)
	return request, nil
}
