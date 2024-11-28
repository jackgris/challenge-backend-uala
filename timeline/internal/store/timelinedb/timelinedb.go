package timelinedb

import (
	"github.com/jackgris/twitter-backend/timeline/internal/domain/timelinemodel"
	"github.com/jackgris/twitter-backend/timeline/internal/store"
)

type Store struct {
	db store.PgxIface
}

func NewStore(db store.PgxIface) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetTimeline(userID string) ([]timelinemodel.Tweet, error) {
	return []timelinemodel.Tweet{}, nil
}

func (s *Store) UpdateTimeline(userID, tweetID string) ([]timelinemodel.Tweet, error) {
	return []timelinemodel.Tweet{}, nil
}
