package main

import (
	"context"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/davecgh/go-spew/spew"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	if os.Getenv("AWS_SAM_LOCAL") != "" {
	}
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, param interface{}) ([]Tweet, error) {
	log.Println(param)
	config := loadConfig()
	oauthConfig := oauth1.NewConfig(config.TwitterConsumerKey, config.TwitterConsumerSecret)
	token := oauth1.NewToken(config.TwitterAccessToken, config.TwitterAccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	trueValue := true
	falseValue := false
	userTimelineParams := &twitter.UserTimelineParams{
		SinceID:         0,
		ScreenName:      "dekokun_test",
		Count:           10,
		IncludeRetweets: &falseValue,
		ExcludeReplies:  &falseValue,
		TrimUser:        &trueValue,
	}
	tweets, responce, err := client.Timelines.UserTimeline(userTimelineParams)
	if err != nil {
		log.Println("responce:", responce)
		log.Println("err:", err)
		return nil, err
	}
	result := []Tweet{}
	for _, tweet := range tweets {
		var mediaURL string
		var expandedURL string
		if tweet.Entities.Media == nil {
			mediaURL = ""
			expandedURL = ""
		} else {
			mediaURL = tweet.Entities.Media[0].MediaURLHttps
			expandedURL = tweet.Entities.Media[0].ExpandedURL
		}
		spew.Dump(mediaURL)
		spew.Dump(expandedURL)
		result = append(result, Tweet{
			Coordinates: tweet.Coordinates,
			CreatedAt:   tweet.CreatedAt,
			Text:        tweet.Text,
			mediaURL:    mediaURL,
			expandedURL: expandedURL,
		})
	}
	spew.Dump(result)

	return result, nil
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
