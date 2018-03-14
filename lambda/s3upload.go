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
	// S3 Client test
	sess := session.New(&aws.Config{})
	s3Svc := s3.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))
	uploader := s3manager.NewUploaderWithClient(s3Svc)
	contentType := "application/json"
	bucket := os.Getenv("BucketName")
	key := "hogehogeeeeetests"
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := s3Svc.GetObject(input)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var oldTweets []Tweet
	if err := decoder.Decode(&oldTweets); err != nil {
		// handle error
		return "", err
	}

	spew.Dump("result:")
	spew.Dump(oldTweets)
	allTweets := append(tweets, oldTweets...)
	m := make(map[int64]bool)
	uniqTweets := []Tweet{}
	// uniq by tweet id
	for _, tweet := range allTweets {
		if !m[tweet.ID] {
			m[tweet.ID] = true
			uniqTweets = append(uniqTweets, tweet)
		}
	}
	jsonBytes, err := json.Marshal(uniqTweets)
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

	uploadResult, err := uploader.Upload(upParams)
	if err != nil {
		log.Println(uploadResult)
		log.Println(err)
		return "", err
	}
	return "s3upload.go", nil
}
