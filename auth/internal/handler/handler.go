package handler

import (
	"net/http"

	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
	"github.com/jackgris/twitter-backend/auth/pkg/logger"
	"github.com/jackgris/twitter-backend/auth/pkg/middleware"
	"github.com/jackgris/twitter-backend/auth/pkg/msgbroker"
)

type UserHandler struct {
	logs      *logger.Logger
	store     Store
	msgBroker *msgbroker.MsgBroker
}

func NewTweetHandler(store Store, msgBroker *msgbroker.MsgBroker, logs *logger.Logger) UserHandler {
	return UserHandler{
		store:     store,
		logs:      logs,
		msgBroker: msgBroker,
	}
}

func NewHandler(store Store, msgBroker *msgbroker.MsgBroker, logs *logger.Logger) (*http.ServeMux, *UserHandler) {
	u := UserHandler{
		store:     store,
		logs:      logs,
		msgBroker: msgBroker,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helthz", middleware.LogResponse(healthCheckHandler, u.logs))
	mux.HandleFunc("POST /create", middleware.LogResponse(u.CreateUser, u.logs))
	mux.HandleFunc("GET /id/{id}", middleware.LogResponse(u.GetUserbyID, u.logs))
	mux.HandleFunc("GET /name/{name}", middleware.LogResponse(u.GetUserbyUsername, u.logs))
	mux.HandleFunc("DELETE /delete/{id}", middleware.LogResponse(u.Delete, u.logs))
	mux.HandleFunc("POST /follow", middleware.LogResponse(u.Follow, u.logs))
	mux.HandleFunc("DELETE /unfollow", middleware.LogResponse(u.Unfollow, u.logs))
	mux.HandleFunc("PATCH /update", middleware.LogResponse(u.Update, u.logs))

	return mux, &u
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Store interface {
	Create(user usermodel.User) (usermodel.User, error)
	GetUserbyID(id string) (usermodel.User, error)
	GetUserbyUsername(username string) (usermodel.User, error)
	Delete(id string) error
	Follow(follow usermodel.UserFollowers) error
	Unfollow(follow usermodel.UserFollowers) error
	Update(user usermodel.User) (usermodel.User, error)
}
