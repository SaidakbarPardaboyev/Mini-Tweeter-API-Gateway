package models

type TweetMedia struct {
	Id        string `json:"id"`
	TweetId   string `json:"tweet_id"`
	MediaType string `json:"media_type"`
	Url       string `json:"url"`
}

type CreateTweetMedia struct {
	TweetId   string `json:"tweet_id"`
	MediaType string `json:"media_type"`
	Url       string `json:"url"`
}

type UpdateTweetMedia struct {
	Id        string `json:"id"`
	TweetId   string `json:"tweet_id"`
	MediaType string `json:"media_type"`
	Url       string `json:"url"`
}
