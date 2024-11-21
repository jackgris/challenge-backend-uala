package tweetdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackgris/twitter-backend/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/twitter-backend/tweet/internal/store/tweetdb"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetLikesByTweetID(t *testing.T) {
	mockLikes := [][]any{
		{
			"csvqda265b6s73dtmot0", // id
			"csvr2keek44s73e2qf90", // tweet_id
			"csvqvamek44s73e2qf8g", // user_id
		},
		{
			"csvqvamek44s73e2qf8g",
			"csvr2omek44s73e2qf9g",
			"csvqda265b6s73dtmot0",
		},
	}

	t.Run("Get Likes OK", func(t *testing.T) {
		mock, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		defer mock.Close(ctx)

		rows := pgxmock.NewRows([]string{"id", "tweet_id", "user_id"}).AddRows(mockLikes...)

		mock.ExpectQuery("SELECT").WithArgs("csvqvamek44s73e2qf8g").WillReturnRows(rows)

		store := tweetdb.NewStore(mock)

		likes, err := store.GetLikesByTweetID(ctx, "csvqvamek44s73e2qf8g")
		assert.NoError(t, err)
		assert.NotNil(t, likes)
		assert.Equal(t, 2, len(likes))

	})

	t.Run("Get database error", func(t *testing.T) {
		mock, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.Background()

		defer mock.Close(ctx)

		rows := pgxmock.NewRows([]string{"id", "tweet_id", "user_id"}).AddRows(mockLikes...).RowError(0, fmt.Errorf("query execution failed"))

		mock.ExpectQuery("SELECT").WithArgs("csvqvamek44s73e2qf8g").WillReturnRows(rows)

		store := tweetdb.NewStore(mock)

		likes, err := store.GetLikesByTweetID(ctx, "csvqvamek44s73e2qf8g")
		assert.Error(t, err)
		assert.Nil(t, likes)
		assert.EqualError(t, err, "row scanning failed: query execution failed")
	})
}

func TestGetTweetByID(t *testing.T) {

	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	store := tweetdb.NewStore(mock)

	tweetID := "csvqda265b6s73dtmot0"
	userID1 := "csvr2keek44s73e2af90"
	userID2 := "csvqda265b6s73dtmot0"
	expectedTweet := tweetmodel.Tweet{
		Id:           tweetID,
		UserID:       userID1,
		Content:      "Test tweet",
		CreatedAt:    time.Now(),
		Encoded_date: "2024-01-01",
		LikeCount:    2,
		RetweetCount: 2,
	}

	mock.ExpectQuery("SELECT id, user_id, content, created_at, encoded_date, like_count, retweet_count").
		WithArgs(tweetID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "content", "created_at", "encoded_date", "like_count", "retweet_count"}).
			AddRow(expectedTweet.Id, expectedTweet.UserID, expectedTweet.Content, expectedTweet.CreatedAt, expectedTweet.Encoded_date, expectedTweet.LikeCount, expectedTweet.RetweetCount))

	mock.ExpectQuery("SELECT id, tweet_id, user_id FROM likes").
		WithArgs(tweetID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "tweet_id", "user_id"}).
			AddRow("csvqda265b6s73dtmot0", tweetID, userID1).
			AddRow("csvr2keek44s73e2af90", tweetID, userID2))

	mock.ExpectQuery("SELECT id, tweet_id, user_id FROM retweets").
		WithArgs(tweetID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "tweet_id", "user_id"}).
			AddRow("csvqda265b6s73dtmot0", tweetID, userID1).
			AddRow("csvr2keek44s73e2af90", tweetID, userID2))

	tweet, err := store.GetByID(tweetID)

	assert.NoError(t, err)
	assert.Equal(t, tweetID, tweet.Id)
	assert.Equal(t, 2, len(tweet.Likes), "Amount of likes")
	assert.Equal(t, 2, len(tweet.Retweets), "Amount of retweets")

	assert.NoError(t, mock.ExpectationsWereMet())
}
