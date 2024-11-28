package handler

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackgris/twitter-backend/tweet/pkg/msgbroker"
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

type FollowersEvent struct {
	Header      msgbroker.Header `json:"header"`
	TweetID     string           `json:"tweet_id"`
	Content     string           `json:"content"`
	FollowersID []string         `json:"followers_id"`
}

type TweetEvent struct {
	Header  msgbroker.Header `json:"header"`
	UserID  string           `json:"user_id"`
	TweetID string           `json:"tweet_id"`
	Content string           `json:"content"`
}

func NewTweet(tweetID, userID, content string) *message.Message {
	event := TweetEvent{
		Header:  msgbroker.NewHeader("tweet"),
		TweetID: tweetID,
		Content: content,
		UserID:  userID,
	}
	tweetMsg, _ := json.Marshal(event)

	return message.NewMessage(event.Header.ID, tweetMsg)
}

func (t TweetHandler) SendTweetToFollowersEvent() {
	ctx := context.Background()
	messages, err := t.msgBroker.SubscribeEvents("followers")
	if err != nil {
		t.logs.Error(ctx, "tweet service", "reading paylod followers", err)
		return
	}

	for msg := range messages {
		msg.Ack()
		followers := FollowersEvent{}
		err := json.Unmarshal(msg.Payload, &followers)
		if err != nil {
			t.logs.Error(ctx, "tweet service", "reading paylod followers", err)
			continue
		}

		t.logs.Info(ctx, "tweet service", "publish message", "topic", "followers", followers, "msg ID", msg.UUID)
		for _, follower := range followers.FollowersID {
			if follower != "" {
				tweet := NewTweet(followers.TweetID, follower, followers.Content)
				t.msgBroker.PublishMessages("tweet-"+follower, tweet)
			}
		}
	}
}
