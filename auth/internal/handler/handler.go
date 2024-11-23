package handler

import (
	"net/http"

	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
	"github.com/jackgris/twitter-backend/auth/pkg/logger"
	"github.com/jackgris/twitter-backend/auth/pkg/middleware"
)

type UserHandler struct {
	logs  *logger.Logger
	store Store
}

func NewTweetHandler(store Store, logs *logger.Logger) UserHandler {
	return UserHandler{
		store: store,
		logs:  logs,
	}
}

func NewHandler(store Store, logs *logger.Logger) *http.ServeMux {
	u := UserHandler{
		store: store,
		logs:  logs,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /helthz", middleware.LogResponse(healthCheckHandler, u.logs))
	mux.HandleFunc("POST /create", middleware.LogResponse(u.CreateUser, u.logs))
	mux.HandleFunc("GET /id/{id}", middleware.LogResponse(u.GetUserbyID, u.logs))
	mux.HandleFunc("GET /name/{name}", middleware.LogResponse(u.GetUserbyUsername, u.logs))
	mux.HandleFunc("DELETE /delete/{id}", middleware.LogResponse(u.Delete, u.logs))
	mux.HandleFunc("POST /like", middleware.LogResponse(u.Follow, u.logs))
	mux.HandleFunc("DELETE /like", middleware.LogResponse(u.Unfollow, u.logs))
	mux.HandleFunc("PUT /retweet", middleware.LogResponse(u.Update, u.logs))

	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

type Store interface {
	Create(user usermodel.User) (usermodel.User, error)
	GetUserbyID(id string) (*usermodel.User, error)
	GetUserbyUsername(username string) (*usermodel.User, error)
	Delete(id string) error
	Follow(follow usermodel.UserFollowers) error
	Unfollow(follow usermodel.UserFollowers) error
	Update(user usermodel.User) (usermodel.User, error)
}
