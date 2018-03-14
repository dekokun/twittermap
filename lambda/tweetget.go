package main

import (
	"context"
	"log"
	"os"
	"strconv"

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
	screenName := os.Getenv("TwitterUseName")
	userTimelineParams := &twitter.UserTimelineParams{
		ScreenName:      screenName,
		Count:           20,
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
		if tweet.Place == nil {
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
			Coordinates: tweet.Place.BoundingBox.Coordinates[0][0],
			CreatedAt:   tweet.CreatedAt,
			Text:        tweet.Text,
			Url:         "https://twitter.com/" + screenName + "/status/" + strconv.FormatInt(tweet.ID, 10),
			MediaURL:    mediaURL,
			ExpandedURL: expandedURL,
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
