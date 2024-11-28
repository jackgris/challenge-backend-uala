package timelinedb

import (
	"time"

	"github.com/jackgris/twitter-backend/timeline/internal/domain/timelinemodel"
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

func TweetToModel(tweet Tweet) timelinemodel.Tweet {
	likes := []timelinemodel.Like{}
	for _, like := range tweet.Likes {
		newLike := timelinemodel.Like{
			Id:      like.Id,
			TweetID: like.TweetID,
			UserID:  like.UserID,
		}
		likes = append(likes, newLike)
	}

	retweets := []timelinemodel.Retweet{}
	for _, retweet := range tweet.Retweets {
		newRetweet := timelinemodel.Retweet{
			Id:      retweet.Id,
			TweetID: retweet.TweetID,
			UserID:  retweet.UserID,
		}
		retweets = append(retweets, newRetweet)
	}

	return timelinemodel.Tweet{
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

func RetweetToModel(retweet Retweet) timelinemodel.Retweet {
	return timelinemodel.Retweet{
		TweetID: retweet.TweetID,
		UserID:  retweet.UserID,
	}
}

func LikeToModel(like Like) timelinemodel.Like {
	return timelinemodel.Like{
		TweetID: like.TweetID,
		UserID:  like.UserID,
	}
}
