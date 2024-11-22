package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/twitter-backend/tweet/pkg/uuid"
)

func (t TweetHandler) LikeTweet(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TweetID string `json:"tweet_id"`
		UserID  string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if input.TweetID == "" || input.UserID == "" {
		http.Error(w, "tweet_id and user_id are required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.UserID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.TweetID); !ok {
		http.Error(w, "tweet id invalid", http.StatusBadRequest)
		return
	}

	like := tweetmodel.Like{
		TweetID: input.TweetID,
		UserID:  input.UserID,
	}

	like, err := t.store.Like(like)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add like to tweet ID: %s", like.TweetID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LikeToJSON(like))
}

func (t TweetHandler) DislikeTweet(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ID      string `json:"id"`
		TweetID string `json:"tweet_id"`
		UserID  string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if input.TweetID == "" || input.UserID == "" {
		http.Error(w, "tweet_id and user_id are required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.UserID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.TweetID); !ok {
		http.Error(w, "tweet id invalid", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.ID); !ok {
		http.Error(w, "like id invalid", http.StatusBadRequest)
		return
	}

	like := tweetmodel.Like{
		Id:      input.ID,
		TweetID: input.TweetID,
		UserID:  input.UserID,
	}

	err := t.store.Dislike(like)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove like from tweet ID: %s", like.Id), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Tweet dislike successful"))
}
