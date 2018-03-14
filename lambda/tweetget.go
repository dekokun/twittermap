package main

import (
	"context"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/lambda"
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
		ScreenName:      os.Getenv("TwitterUseName"),
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
		if tweet.Coordinates == nil {
			continue
		}
		var mediaURL string
		var expandedURL string
		if tweet.Entities.Media == nil {
			mediaURL = ""
			expandedURL = ""
		} else {
			mediaURL = tweet.Entities.Media[0].MediaURLHttps
			expandedURL = tweet.Entities.Media[0].ExpandedURL
		}
		result = append(result, Tweet{
			ID:          tweet.ID,
			Coordinates: tweet.Coordinates,
			CreatedAt:   tweet.CreatedAt,
			Text:        tweet.Text,
			mediaURL:    mediaURL,
			expandedURL: expandedURL,
		})
	}

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
