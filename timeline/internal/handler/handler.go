package handler

import (
	"net/http"

	"github.com/jackgris/twitter-backend/timeline/internal/domain/timelinemodel"
	"github.com/jackgris/twitter-backend/timeline/pkg/logger"
	"github.com/jackgris/twitter-backend/timeline/pkg/middleware"
	"github.com/jackgris/twitter-backend/timeline/pkg/msgbroker"
)

type TimelineHandler struct {
	logs      *logger.Logger
	store     Store
	msgBroker *msgbroker.MsgBroker
}

func NewTweetHandler(store Store, msgBroker *msgbroker.MsgBroker, logs *logger.Logger) TimelineHandler {
	return TimelineHandler{
		store:     store,
		msgBroker: msgBroker,
		logs:      logs,
	}
}

func NewHandler(store Store, msgBroker *msgbroker.MsgBroker, logs *logger.Logger) *http.ServeMux {
	t := TimelineHandler{
		store:     store,
		msgBroker: msgBroker,
		logs:      logs,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helthz", middleware.LogResponse(healthCheckHandler, t.logs))
	mux.HandleFunc("GET /timeline", middleware.LogResponse(t.GetTimelineHandler, t.logs))
	mux.HandleFunc("GET /update", middleware.LogResponse(t.UpdateTimelineHandler, t.logs))

	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Store interface {
	GetTimeline(userID string) ([]timelinemodel.Tweet, error)
	UpdateTimeline(userID, tweetID string) ([]timelinemodel.Tweet, error)
}

func (t *TimelineHandler) GetTimelineHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (t *TimelineHandler) UpdateTimelineHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
