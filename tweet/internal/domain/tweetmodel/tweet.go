package tweetmodel

import "time"

type Tweet struct {
	Id           string
	UserID       string
	Content      string
	CreatedAt    time.Time
	Encoded_date string
	LikeCount    int
	RetweetCount int
	Likes        []Like
	Retweets     []Retweet
}

type Retweet struct {
	Id      string
	TweetID string
	UserID  string
}

type Like struct {
	Id      string
	TweetID string
	UserID  string
}
