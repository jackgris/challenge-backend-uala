package userdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
	"github.com/jackgris/twitter-backend/auth/internal/store"
	"github.com/jackgris/twitter-backend/auth/pkg/uuid"
)

type Store struct {
	db store.PgxIface
}

func NewStore(db store.PgxIface) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Create(user usermodel.User) (usermodel.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()
	query := `
                  INSERT INTO users (id, username, email, password, follower_count, following_count, salt, token, date_created, encoded_date)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                  RETURNING id, username, email, password, follower_count, following_count, salt, token, date_created, encoded_date;
        `
	var newUser usermodel.User
	err := s.db.QueryRow(ctx, query,
		user.ID,
		user.UserName,
		user.Email,
		user.Password,
		user.FollowerCount,
		user.FollowingCount,
		user.Salt,
		user.Token,
		user.DateCreated,
		user.EncodedDate,
	).Scan(&newUser.ID,
		&newUser.UserName,
		&newUser.Email,
		&newUser.Password,
		&newUser.FollowerCount,
		&newUser.FollowingCount,
		&newUser.Salt,
		&newUser.Token,
		&newUser.DateCreated,
		&newUser.EncodedDate)
	if err != nil {
		return usermodel.User{}, fmt.Errorf("failed to insert tweet: %w", err)
	}

	return newUser, nil
}

var ErrDeleteUser = errors.New("no user found")

func (s *Store) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
		DELETE FROM users
		WHERE id = $1;
	`

	commandTag, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tweet: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.Join(ErrDeleteUser, fmt.Errorf("with ID: %s", id))
	}

	return nil
}

func (s *Store) Follow(follow usermodel.UserFollowers) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
                INSERT INTO user_followers (id, user_id, username) VALUES ($1, $2, $3)
        `
	_, err := s.db.Exec(ctx, query,
		uuid.New(),
		follow.UserID,
		follow.UserName,
	)
	if err != nil {
		return err
	}

	query = `
               UPDATE users SET follower_count = follower_count + 1 WHERE username = $1
        `
	_, err = s.db.Exec(ctx, query, follow.UserName)
	if err != nil {
		return err
	}

	query = `
               UPDATE users SET following_count = following_count + 1 WHERE username = $1
        `
	_, err = s.db.Exec(ctx, query, follow.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Unfollow(follow usermodel.UserFollowers) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
                DELETE FROM user_followers WHERE user_id = $1 AND username = $2
        `
	_, err := s.db.Exec(ctx, query,
		follow.UserID,
		follow.UserName,
	)
	if err != nil {
		return err
	}

	query = `
               UPDATE users SET follower_count = follower_count - 1 WHERE username = $1
        `
	_, err = s.db.Exec(ctx, query, follow.UserName)
	if err != nil {
		return err
	}

	query = `
               UPDATE users SET following_count = following_count - 1 WHERE username = $1
        `
	_, err = s.db.Exec(ctx, query, follow.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserbyUsername(username string) (*usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
                SELECT * FROM users WHERE username = $1
        `
	var user usermodel.User
	err := s.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.FollowerCount,
		&user.FollowingCount,
		&user.Salt, &user.Token,
		&user.DateCreated,
		&user.EncodedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserbyID(id string) (*usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
                SELECT * FROM users WHERE id = $1
        `
	var user usermodel.User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.FollowerCount,
		&user.FollowingCount,
		&user.Salt,
		&user.Token,
		&user.DateCreated,
		&user.EncodedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) Update(user usermodel.User) (usermodel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
	defer cancel()

	query := `
                UPDATE users
                SET username = $1, email = $2, password = $3, follower_count = $4,
                    following_count = $5, salt = $6, token = $7, date_created = $8,
                    encoded_date = $9 WHERE id = $10
        `
	_, err := s.db.Exec(ctx, query,
		user.UserName,
		user.Email,
		user.Password,
		user.FollowerCount,
		user.FollowingCount,
		user.Salt,
		user.Token,
		user.DateCreated,
		user.EncodedDate,
		user.ID,
	)
	if err != nil {
		return usermodel.User{}, err
	}

	return user, nil
}
