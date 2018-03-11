package main

import (
	"github.com/dghubble/go-twitter/twitter"
)

type Tweet struct {
	ID          int64                `json:"id"`
	Coordinates *twitter.Coordinates `json:"coordinates"`
	CreatedAt   string               `json:"created_at"`
	Text        string               `json:"text"`
	mediaURL    string               `json:"media_url"`
	expandedURL string               `json:"expanded_url"`
}
