package handler

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackgris/twitter-backend/timeline/pkg/msgbroker"
)

type TweetCreated struct {
	Header  msgbroker.Header `json:"header"`
	UserID  string           `json:"user_id"`
	TweetID string           `json:"tweet_id"`
	Content string           `json:"content"`
}

func NewTweetCreated(userID, tweetID, content string) *message.Message {
	event := TweetCreated{
		Header:  msgbroker.NewHeader("tweet_created"),
		UserID:  userID,
		TweetID: tweetID,
		Content: content,
	}
	tweetMsg, _ := json.Marshal(event)

	return message.NewMessage(event.Header.ID, tweetMsg)
}

type GetFollowers struct {
	Header msgbroker.Header `json:"header"`
	UserID string           `json:"user_id"`
}

func NewGetFollowers(userID string) *message.Message {
	event := GetFollowers{
		Header: msgbroker.NewHeader("get_followers"),
		UserID: userID,
	}
	getFollowersMsg, _ := json.Marshal(event)

	return message.NewMessage(event.Header.ID, getFollowersMsg)
}
