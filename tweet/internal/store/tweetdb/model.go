package tweetdb

import (
	"time"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
)

type Tweet struct {
	Id           string
	UserID       string
	Content      string
	CreatedAt    time.Time
	EncodedDate  string
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

func TweetToModel(tweet Tweet) tweetmodel.Tweet {
	likes := []tweetmodel.Like{}
	for _, like := range tweet.Likes {
		newLike := tweetmodel.Like{
			Id:      like.Id,
			TweetID: like.TweetID,
			UserID:  like.UserID,
		}
		likes = append(likes, newLike)
	}

	retweets := []tweetmodel.Retweet{}
	for _, retweet := range tweet.Retweets {
		newRetweet := tweetmodel.Retweet{
			Id:      retweet.Id,
			TweetID: retweet.TweetID,
			UserID:  retweet.UserID,
		}
		retweets = append(retweets, newRetweet)
	}

	return tweetmodel.Tweet{
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

func RetweetToModel(retweet Retweet) tweetmodel.Retweet {
	return tweetmodel.Retweet{
		TweetID: retweet.TweetID,
		UserID:  retweet.UserID,
	}
}

func LikeToModel(like Like) tweetmodel.Like {
	return tweetmodel.Like{
		TweetID: like.TweetID,
		UserID:  like.UserID,
	}
}
