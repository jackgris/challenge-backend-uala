package handler

import (
	"net/http"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
	"github.com/jackgris/twitter-backend/tweet/pkg/middleware"
	"github.com/jackgris/twitter-backend/tweet/pkg/msgbroker"
)

type TweetHandler struct {
	logs      *logger.Logger
	store     Store
	msgBroker *msgbroker.MsgBroker
}

func NewTweetHandler(store Store, msgBroker *msgbroker.MsgBroker, logs *logger.Logger) TweetHandler {
	return TweetHandler{
		store:     store,
		msgBroker: msgBroker,
		logs:      logs,
	}
}

func NewHandler(store Store, msgBroker *msgbroker.MsgBroker, logs *logger.Logger) (*http.ServeMux, *TweetHandler) {
	t := TweetHandler{
		store:     store,
		msgBroker: msgBroker,
		logs:      logs,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helthz", middleware.LogResponse(healthCheckHandler, t.logs))
	mux.HandleFunc("GET /id/{id}", middleware.LogResponse(t.GetTweetById, t.logs))
	mux.HandleFunc("POST /create", middleware.LogResponse(t.CreateTweet, t.logs))
	mux.HandleFunc("DELETE /delete/{id}", middleware.LogResponse(t.DeleteTweet, t.logs))
	mux.HandleFunc("POST /like", middleware.LogResponse(t.LikeTweet, t.logs))
	mux.HandleFunc("DELETE /like", middleware.LogResponse(t.DislikeTweet, t.logs))
	mux.HandleFunc("POST /retweet", middleware.LogResponse(t.ReTweet, t.logs))
	mux.HandleFunc("DELETE /retweet", middleware.LogResponse(t.DeleteReTweet, t.logs))

	return mux, &t
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Store interface {
	GetByID(id string) (tweetmodel.Tweet, error)
	GetByUser(userID string) ([]tweetmodel.Tweet, error)
	Create(t tweetmodel.Tweet) (tweetmodel.Tweet, error)
	Delete(tweetID string) error
	Like(like tweetmodel.Like) (tweetmodel.Like, error)
	Dislike(like tweetmodel.Like) error
	ReTweet(retweet tweetmodel.Retweet) (tweetmodel.Retweet, error)
	DeleteReTweet(retweet tweetmodel.Retweet) error
}
