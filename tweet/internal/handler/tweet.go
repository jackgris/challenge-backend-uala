package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/twitter-backend/tweet/internal/store/tweetdb"
	"github.com/jackgris/twitter-backend/tweet/pkg/uuid"
)

func (t TweetHandler) createTweet(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if input.UserID == "" || input.Content == "" {
		http.Error(w, "user_id and content are required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.UserID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(input.Content) > 280 {
		http.Error(w, "content size should have a maximun of 280 characters", http.StatusBadRequest)
		return
	}

	tweet := tweetmodel.Tweet{
		UserID:  input.UserID,
		Content: input.Content,
	}
	tweet, err := t.store.Create(tweet)
	if err != nil {
		http.Error(w, "Can't save tweet in database", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(TweetToJSON(tweet))

}

func (t TweetHandler) getTweetById(w http.ResponseWriter, r *http.Request) {

	tweetID := r.PathValue("id")
	if tweetID == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(tweetID); !ok {
		http.Error(w, "tweet id invalid", http.StatusBadRequest)
		return
	}

	tweet, err := t.store.GetByID(tweetID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Tweet not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to retrieve tweet: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(TweetToJSON(tweet))
}

func (t TweetHandler) deleteTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := r.PathValue("id")
	if tweetID == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(tweetID); !ok {
		http.Error(w, "tweet id invalid", http.StatusBadRequest)
		return
	}

	err := t.store.Delete(tweetID)
	if errors.Is(err, tweetdb.ErrDeleteTweet) {
		http.Error(w, "Tweet not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
