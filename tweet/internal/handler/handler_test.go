package handler_test

import (
	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
)

type MockStore struct {
	LikeFunc func(like tweetmodel.Like) (tweetmodel.Like, error)
}

func (m *MockStore) GetByID(id string) (tweetmodel.Tweet, error) {
	return tweetmodel.Tweet{}, nil
}
func (m *MockStore) GetByUser(userID string) ([]tweetmodel.Tweet, error) {
	return nil, nil
}
func (m *MockStore) Create(t tweetmodel.Tweet) (tweetmodel.Tweet, error) {
	return tweetmodel.Tweet{}, nil
}
func (m *MockStore) Delete(tweetID string) error {
	return nil
}
func (m *MockStore) Like(like tweetmodel.Like) (tweetmodel.Like, error) {
	return m.LikeFunc(like)
}
func (m *MockStore) Dislike(like tweetmodel.Like) error {
	return nil
}
func (m *MockStore) ReTweet(retweet tweetmodel.Retweet) (tweetmodel.Retweet, error) {
	return tweetmodel.Retweet{}, nil
}
func (m *MockStore) DeleteReTweet(retweet tweetmodel.Retweet) error {
	return nil
}
