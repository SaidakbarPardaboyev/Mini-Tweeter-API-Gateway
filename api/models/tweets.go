package models

type Tweet struct {
	Id        string        `json:"id"`
	Content   string        `json:"content"`
	UserId    string        `json:"user_id"`
	Medias    []*TweetMedia `json:"medias"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

type CreateTweet struct {
	Content string `json:"content"`
	UserId  string `json:"user_id"`
}

type UpdateTweet struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	UserId  string `json:"user_id"`
}
