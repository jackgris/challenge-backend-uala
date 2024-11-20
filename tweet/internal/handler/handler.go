package handler

import (
	"net/http"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
	"github.com/jackgris/twitter-backend/tweet/pkg/middleware"
)

type TweetHandler struct {
	logs  *logger.Logger
	store Store
}

func NewHandler(store Store, logs *logger.Logger) *http.ServeMux {
	t := TweetHandler{
		store: store,
		logs:  logs,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helthz", middleware.LogResponse(healthCheckHandler, t.logs))
	mux.HandleFunc("GET /id", middleware.LogResponse(t.getTweetById, t.logs))
	mux.HandleFunc("POST /create", middleware.LogResponse(t.createTweet, t.logs))
	mux.HandleFunc("DELETE /id/{id}/delete", middleware.LogResponse(t.deleteTweet, t.logs))
	mux.HandleFunc("POST /id/{id}/like", middleware.LogResponse(t.likeTweet, t.logs))
	mux.HandleFunc("DELETE /id/{id}/dislike", middleware.LogResponse(t.dislikeTweet, t.logs))
	mux.HandleFunc("POST /id/{id}/retweet", middleware.LogResponse(t.reTweet, t.logs))
	mux.HandleFunc("DELETE /id/{id}/retweet", middleware.LogResponse(t.deleteReTweet, t.logs))

	return mux
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
