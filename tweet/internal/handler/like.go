package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackgris/challenge-backend-uala/tweet/internal/domain/tweetmodel"
	"github.com/rs/xid"
)

func (t TweetHandler) likeTweet(w http.ResponseWriter, r *http.Request) {
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

	id, err := xid.FromString(input.TweetID)
	if err != nil {
		http.Error(w, "tweet id invalid", http.StatusBadRequest)
		return
	}

	like := tweetmodel.Like{
		TweetID: id.String(),
		UserID:  input.UserID,
	}

	like, err = t.store.Like(like)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add like to tweet ID: %s", like.TweetID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LikeToJSON(like))
}

func (t TweetHandler) dislikeTweet(w http.ResponseWriter, r *http.Request) {

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

	tweetID, err := xid.FromString(input.TweetID)
	if err != nil {
		http.Error(w, "tweet id invalid", http.StatusBadRequest)
		return
	}

	likeID, err := xid.FromString(input.ID)
	if err != nil {
		http.Error(w, "like id invalid", http.StatusBadRequest)
		return
	}

	like := tweetmodel.Like{
		Id:      likeID.String(),
		TweetID: tweetID.String(),
		UserID:  input.UserID,
	}

	err = t.store.Dislike(like)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove like from tweet ID: %s", like.Id), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Tweet dislike successful"))
}
