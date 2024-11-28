package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/twitter-backend/tweet/internal/handler"
	"github.com/jackgris/twitter-backend/tweet/pkg/logger"
	"github.com/jackgris/twitter-backend/tweet/pkg/msgbroker"
	"github.com/jackgris/twitter-backend/tweet/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLikeTweet(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		mockLikeFunc func(like tweetmodel.Like) (tweetmodel.Like, error)
		expectedCode int
		expectedBody string
	}{
		{
			name: "Success",
			requestBody: map[string]string{
				"tweet_id": uuid.New(),
				"user_id":  uuid.New(),
			},
			mockLikeFunc: func(like tweetmodel.Like) (tweetmodel.Like, error) {
				return like, nil
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"tweet_id":`,
		},
		{
			name:         "Invalid JSON payload",
			requestBody:  `invalid json`,
			mockLikeFunc: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid JSON payload",
		},
		{
			name: "Missing fields",
			requestBody: map[string]string{
				"tweet_id": "",
				"user_id":  "",
			},
			mockLikeFunc: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "tweet_id and user_id are required",
		},
		{
			name: "Invalid UserID",
			requestBody: map[string]string{
				"tweet_id": uuid.New(),
				"user_id":  "invalid-uuid",
			},
			mockLikeFunc: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "user id invalid",
		},
		{
			name: "Store Like error",
			requestBody: map[string]string{
				"tweet_id": uuid.New(),
				"user_id":  uuid.New(),
			},
			mockLikeFunc: func(like tweetmodel.Like) (tweetmodel.Like, error) {
				return tweetmodel.Like{}, errors.New("store error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Failed to add like to tweet ID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStore := &MockStore{
				LikeFunc: test.mockLikeFunc,
			}

			log := logger.New(io.Discard)
			msgbroker := msgbroker.NewMockMsgBroker(log)
			handler := handler.NewTweetHandler(mockStore, msgbroker, log)

			body, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/like", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.LikeTweet(rec, req)

			assert.Equal(t, test.expectedCode, rec.Code)

			if test.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), test.expectedBody)
			}
		})
	}
}
