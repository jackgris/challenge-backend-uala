package handler

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackgris/twitter-backend/auth/pkg/msgbroker"
)

type Tweet struct {
	Header  msgbroker.Header `json:"header"`
	UserID  string           `json:"user_id"`
	TweetID string           `json:"tweet_id"`
	Content string           `json:"content"`
}

type Followers struct {
	Header      msgbroker.Header `json:"header"`
	TweetID     string           `json:"tweet_id"`
	Content     string           `json:"content"`
	FollowersID []string         `json:"followers_id"`
}

func NewFollowers(tweetID, content string, followers []string) *message.Message {
	event := Followers{
		Header:      msgbroker.NewHeader("followers"),
		TweetID:     tweetID,
		Content:     content,
		FollowersID: followers,
	}
	tweetMsg, _ := json.Marshal(event)

	return message.NewMessage(event.Header.ID, tweetMsg)
}

func (u *UserHandler) SubscribeGetFollowers() {
	ctx := context.Background()
	topic := "tweets"
	messages, err := u.msgBroker.SubscribeEvents(topic)
	if err != nil {
		u.logs.Error(ctx, "auth service", "subscriber: can't subscribe "+topic, err)
	}

	for msg := range messages {
		msg.Ack()
		tweet := Tweet{}
		err := json.Unmarshal(msg.Payload, &tweet)
		if err != nil {
			u.logs.Error(ctx, "auth service", "reading paylod "+topic, err, "msg", msg)
			continue
		}

		u.logs.Info(ctx, "auth service", "receive tweets ", tweet, "ID", msg.UUID)

		user, err := u.store.GetUserbyID(tweet.UserID)
		if err != nil {
			if err != pgx.ErrNoRows {
				u.logs.Error(ctx, "auth service", "getting followers "+topic, err)
			}
			continue
		}

		userFollowers := []string{}
		for _, f := range user.Followers {
			userFollowers = append(userFollowers, f.FollowerID)
		}

		followers := NewFollowers(tweet.TweetID, tweet.Content, userFollowers)

		u.msgBroker.PublishMessages("followers", followers)

		newFollow := Followers{}
		_ = json.Unmarshal(followers.Payload, &newFollow)
		u.logs.Info(ctx, "auth service", "send followers ", newFollow)
	}
}
