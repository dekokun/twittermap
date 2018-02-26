package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	oauthConfig := oauth1.NewConfig(config.twitterConsumerKey, config.twitterConsumerSecret)
	token := oauth1.NewToken(config.twitterAccessToken, config.twitterAccessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	userTimelineParams := &twitter.UserTimelineParams{ScreenName: "dekokun", Count: 2}
	tweets, _, _ := client.Timelines.UserTimeline(userTimelineParams)
	fmt.Printf("USER TIMELINE:\n%+v\n", tweets)

	return request, nil
}

type config struct {
	twitterConsumerKey    string
	twitterConsumerSecret string
	twitterAccessToken    string
	twitterAccessSecret   string
}

func loadConfig() config {
	var configToml config
	_, err := toml.DecodeFile("config.toml", &configToml)
	if err != nil {
		panic(err)
	}
	return configToml
}
