package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackgris/challenge-backend-uala/tweet/internal/domain/tweetmodel"
	"github.com/rs/xid"
)

func (t TweetHandler) reTweet(w http.ResponseWriter, r *http.Request) {
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

	retweet := tweetmodel.Retweet{
		TweetID: input.TweetID,
		UserID:  input.UserID,
	}

	retweet, err := t.store.ReTweet(retweet)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retweet to tweet ID: %s", retweet.TweetID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(RetweetToJSON(retweet))
}

func (t TweetHandler) deleteReTweet(w http.ResponseWriter, r *http.Request) {

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

	retweetID, err := xid.FromString(input.ID)
	if err != nil {
		http.Error(w, "retweet id invalid", http.StatusBadRequest)
		return
	}

	retweet := tweetmodel.Retweet{
		Id:      retweetID.String(),
		TweetID: tweetID.String(),
		UserID:  input.UserID,
	}

	err = t.store.DeleteReTweet(retweet)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to remove retweet from tweet ID: %s", retweet.TweetID), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ReTweet removed successful"))

}
