package main

import (
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	if os.Getenv("AWS_SAM_LOCAL") != "" {
	}
	lambda.Start(handleRequest)
}

func handleRequest(request events.CloudWatchEvent) (events.CloudWatchEvent, error) {
	config := loadConfig()
	oauthConfig := oauth1.NewConfig(config.TwitterConsumerKey, config.TwitterConsumerSecret)
	token := oauth1.NewToken(config.TwitterAccessToken, config.TwitterAccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	userTimelineParams := &twitter.UserTimelineParams{ScreenName: "dekokun", Count: 2}
	tweets, responce, err := client.Timelines.UserTimeline(userTimelineParams)
	if err != nil {
		log.Println("responce:", responce)
		log.Println("err:", err)
	}
	// spew.Dump(tweets)
	spew.Dump(tweets[0].Coordinates)
	spew.Dump(tweets[0].CreatedAt)
	spew.Dump(tweets[0].Text)
	spew.Dump(tweets[0].Entities.Media[0].MediaURLHttps)
	spew.Dump(tweets[0].Entities.Media[0].ExpandedURL)
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

type config struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

func loadConfig() config {
	var configToml config
	_, err := toml.DecodeFile("config.toml", &configToml)
	if err != nil {
		panic(err)
	}
	return configToml
}
