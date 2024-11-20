package tweetdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackgris/challenge-backend-uala/tweet/internal/domain/tweetmodel"
	"github.com/jackgris/challenge-backend-uala/tweet/internal/store"
	"github.com/rs/xid"
)

type Store struct {
	db store.PgxIface
}

func NewStore(db store.PgxIface) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetByID(id string) (tweetmodel.Tweet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
		SELECT id, user_id, content, created_at, encoded_date, like_count, retweet_count
		FROM tweets
		WHERE id = $1;
	`
	var tweet Tweet
	err := s.db.QueryRow(ctx, query, id).Scan(
		&tweet.Id,
		&tweet.UserID,
		&tweet.Content,
		&tweet.CreatedAt,
		&tweet.EncodedDate,
		&tweet.LikeCount,
		&tweet.RetweetCount,
	)
	if err != nil {
		return tweetmodel.Tweet{}, fmt.Errorf("failed to fetch tweet: %w", err)
	}

	// Fetch likes and retweets for the tweet
	tweet.Likes, err = s.GetLikesByTweetID(ctx, tweet.Id)
	if err != nil {
		return tweetmodel.Tweet{}, fmt.Errorf("failed to fetch likes: %w", err)
	}

	tweet.Retweets, err = s.GetRetweetsByTweetID(ctx, tweet.Id)
	if err != nil {
		return tweetmodel.Tweet{}, fmt.Errorf("failed to fetch retweets: %w", err)
	}
	return TweetToModel(tweet), nil
}

func (s *Store) GetByUser(userID string) ([]tweetmodel.Tweet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
		SELECT id, user_id, content, created_at, encoded_date, like_count, retweet_count
		FROM tweets
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var tweetsDb []Tweet

	for rows.Next() {
		var tweet Tweet
		err := rows.Scan(
			&tweet.Id,
			&tweet.UserID,
			&tweet.Content,
			&tweet.CreatedAt,
			&tweet.EncodedDate,
			&tweet.LikeCount,
			&tweet.RetweetCount,
		)
		if err != nil {
			return nil, fmt.Errorf("row scanning failed: %w", err)
		}

		tweet.Likes, err = s.GetLikesByTweetID(ctx, tweet.Id)
		if err != nil {
			return nil, fmt.Errorf("fetching likes failed: %w", err)
		}

		tweet.Retweets, err = s.GetRetweetsByTweetID(ctx, tweet.Id)
		if err != nil {
			return nil, fmt.Errorf("fetching retweets failed: %w", err)
		}

		tweetsDb = append(tweetsDb, tweet)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", rows.Err())
	}

	tweetModel := []tweetmodel.Tweet{}
	for _, tweet := range tweetsDb {
		tweetModel = append(tweetModel, TweetToModel(tweet))
	}

	return tweetModel, nil
}

func (s *Store) GetLikesByTweetID(ctx context.Context, tweetID string) ([]Like, error) {
	query := `
		SELECT id, tweet_id, user_id
		FROM likes
		WHERE tweet_id = $1;
	`

	rows, err := s.db.Query(ctx, query, tweetID)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var likes []Like
	for rows.Next() {
		var like Like
		err := rows.Scan(&like.Id, &like.TweetID, &like.UserID)
		if err != nil {
			return nil, fmt.Errorf("row scanning failed: %w", err)
		}
		likes = append(likes, like)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", rows.Err())
	}

	return likes, nil
}

func (s *Store) GetRetweetsByTweetID(ctx context.Context, tweetID string) ([]Retweet, error) {
	query := `
		SELECT id, tweet_id, user_id
		FROM retweets
		WHERE tweet_id = $1;
	`

	rows, err := s.db.Query(ctx, query, tweetID)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var retweets []Retweet
	for rows.Next() {
		var retweet Retweet
		err := rows.Scan(&retweet.Id, &retweet.TweetID, &retweet.UserID)
		if err != nil {
			return nil, fmt.Errorf("row scanning failed: %w", err)
		}
		retweets = append(retweets, retweet)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", rows.Err())
	}

	return retweets, nil
}

func (s *Store) Create(t tweetmodel.Tweet) (tweetmodel.Tweet, error) {

	tweetID := xid.New().String()
	createdAt := time.Now()
	encodedDate := createdAt.Format(time.RFC3339)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
		INSERT INTO tweets (id, user_id, content, created_at, encoded_date, like_count, retweet_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, content, created_at, encoded_date, like_count, retweet_count;
	`

	var tweet tweetmodel.Tweet
	err := s.db.QueryRow(ctx, query, tweetID, t.UserID, t.Content, createdAt, encodedDate, 0, 0).Scan(
		&tweet.Id,
		&tweet.UserID,
		&tweet.Content,
		&tweet.CreatedAt,
		&tweet.Encoded_date,
		&tweet.LikeCount,
		&tweet.RetweetCount,
	)
	if err != nil {
		return tweetmodel.Tweet{}, fmt.Errorf("failed to insert tweet: %w", err)
	}

	return tweet, nil
}

var ErrDeleteTweet = errors.New("no tweet found")

func (s *Store) Delete(tweetID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
		DELETE FROM tweets
		WHERE id = $1;
	`

	commandTag, err := s.db.Exec(ctx, query, tweetID)
	if err != nil {
		return fmt.Errorf("failed to delete tweet: %w", err)
	}

	// Check if the tweet was found and deleted
	if commandTag.RowsAffected() == 0 {
		return errors.Join(ErrDeleteTweet, fmt.Errorf("with ID: %s", tweetID))
	}

	return nil
}

func (s *Store) Like(like tweetmodel.Like) (tweetmodel.Like, error) {
	likeID := xid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return tweetmodel.Like{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	insertQuery := `
		INSERT INTO likes (id, tweet_id, user_id)
		VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(ctx, insertQuery, likeID, like.TweetID, like.UserID)
	if err != nil {
		return tweetmodel.Like{}, fmt.Errorf("failed to insert like: %w", err)
	}

	// Increment the like_count in the tweets table
	updateTweetQuery := `
		UPDATE tweets
		SET like_count = like_count + 1
		WHERE id = $1;
	`
	_, err = tx.Exec(ctx, updateTweetQuery, like.TweetID)
	if err != nil {
		return tweetmodel.Like{}, fmt.Errorf("failed to update like count: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return tweetmodel.Like{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tweetmodel.Like{
		Id:      likeID,
		TweetID: like.TweetID,
		UserID:  like.UserID,
	}, nil
}

func (s *Store) Dislike(like tweetmodel.Like) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	// Begin a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Delete the like from the likes table
	deleteLikeQuery := `
		DELETE FROM likes
		WHERE tweet_id = $1 AND user_id = $2;
	`
	commandTag, err := tx.Exec(ctx, deleteLikeQuery, like.TweetID, like.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	// Check if a like was found and deleted
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no like found for tweet %s by user %s", like.TweetID, like.UserID)
	}

	// Decrement the like_count in the tweets table
	updateTweetQuery := `
		UPDATE tweets
		SET like_count = like_count - 1
		WHERE id = $1;
	`
	_, err = tx.Exec(ctx, updateTweetQuery, like.TweetID)
	if err != nil {
		return fmt.Errorf("failed to update like count: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) ReTweet(retweet tweetmodel.Retweet) (tweetmodel.Retweet, error) {

	retweetID := xid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	// Begin a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return tweetmodel.Retweet{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Insert the retweet into the retweets table
	insertRetweetQuery := `
		INSERT INTO retweets (id, tweet_id, user_id)
		VALUES ($1, $2, $3);
	`
	_, err = tx.Exec(ctx, insertRetweetQuery, retweetID, retweet.TweetID, retweet.UserID)
	if err != nil {
		return tweetmodel.Retweet{}, fmt.Errorf("failed to insert retweet: %w", err)
	}

	// Increment the retweet_count in the tweets table
	updateTweetQuery := `
		UPDATE tweets
		SET retweet_count = retweet_count + 1
		WHERE id = $1;
	`
	_, err = tx.Exec(ctx, updateTweetQuery, retweet.TweetID)
	if err != nil {
		return tweetmodel.Retweet{}, fmt.Errorf("failed to update retweet count: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return tweetmodel.Retweet{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the created Retweet
	return tweetmodel.Retweet{
		Id:      retweetID,
		TweetID: retweet.TweetID,
		UserID:  retweet.UserID,
	}, nil
}

func (s *Store) DeleteReTweet(retweet tweetmodel.Retweet) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	// Begin a transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Delete the retweet from the retweets table
	deleteRetweetQuery := `
		DELETE FROM retweets
		WHERE tweet_id = $1 AND user_id = $2;
	`
	commandTag, err := tx.Exec(ctx, deleteRetweetQuery, retweet.TweetID, retweet.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete retweet: %w", err)
	}

	// Check if a retweet was found and deleted
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no retweet found for tweet %s by user %s", retweet.TweetID, retweet.UserID)
	}

	// Decrement the retweet_count in the tweets table
	updateTweetQuery := `
		UPDATE tweets
		SET retweet_count = retweet_count - 1
		WHERE id = $1;
	`
	_, err = tx.Exec(ctx, updateTweetQuery, retweet.TweetID)
	if err != nil {
		return fmt.Errorf("failed to update retweet count: %w", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
