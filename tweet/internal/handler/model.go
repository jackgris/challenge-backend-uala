package handler

import (
	"time"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
)

type Tweet struct {
	Id           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"tweet_content"`
	CreatedAt    time.Time `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	RetweetCount int       `json:"retweet_count"`
	Likes        []Like    `json:"likes"`
	Retweets     []Retweet `json:"retweets"`
}

type Retweet struct {
	Id      string `json:"id"`
	TweetID string `json:"tweet_id"`
	UserID  string `json:"user_id"`
}

type Like struct {
	Id      string `json:"id"`
	TweetID string `json:"tweet_id"`
	UserID  string `json:"user_id"`
}

func TweetToJSON(tweet tweetmodel.Tweet) Tweet {
	likes := []Like{}
	for _, like := range tweet.Likes {
		newLike := Like{
			Id:      like.Id,
			TweetID: like.TweetID,
			UserID:  like.UserID,
		}
		likes = append(likes, newLike)
	}

	retweets := []Retweet{}
	for _, retweet := range tweet.Retweets {
		newRetweet := Retweet{
			Id:      retweet.Id,
			TweetID: retweet.TweetID,
			UserID:  retweet.UserID,
		}
		retweets = append(retweets, newRetweet)
	}

	return Tweet{
		Id:           tweet.Id,
		UserID:       tweet.UserID,
		Content:      tweet.Content,
		CreatedAt:    tweet.CreatedAt,
		LikeCount:    tweet.LikeCount,
		RetweetCount: tweet.RetweetCount,
		Likes:        likes,
		Retweets:     retweets,
	}
}

func RetweetToJSON(retweet tweetmodel.Retweet) Retweet {
	return Retweet{
		TweetID: retweet.TweetID,
		UserID:  retweet.UserID,
	}
}

func LikeToJSON(like tweetmodel.Like) Like {
	return Like{
		Id:      like.Id,
		TweetID: like.TweetID,
		UserID:  like.UserID,
	}
}
