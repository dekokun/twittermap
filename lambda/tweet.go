package main

type Tweet struct {
	ID          int64      `json:"id"`
	Coordinates [2]float64 `json:"coordinates"`
	CreatedAt   string     `json:"created_at"`
	Text        string     `json:"text"`
	mediaURL    string     `json:"media_url"`
	expandedURL string     `json:"expanded_url"`
}
