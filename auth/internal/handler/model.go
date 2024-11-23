package handler

import (
	"time"

	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
)

type User struct {
	ID             string          `json:"id"`
	UserName       string          `json:"username"`
	Email          string          `json:"email"`
	Password       string          `json:"-"`
	FollowerCount  int             `json:"follower_count"`
	FollowingCount int             `json:"following_count"`
	Salt           string          `json:"-"`
	Token          string          `json:"-"`
	DateCreated    time.Time       `json:"date_created"`
	EncodedDate    string          `json:"-"`
	Followers      []UserFollowers `json:"followers"`
	Following      []UserFollowers `json:"following"`
}

type UserFollowers struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
}

func UserToJSON(user usermodel.User) User {
	return User{
		ID:             user.ID,
		UserName:       user.UserName,
		Email:          user.Email,
		Password:       user.Password,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		Salt:           user.Salt,
		Token:          user.Token,
		DateCreated:    user.DateCreated,
		EncodedDate:    user.EncodedDate,
		Followers:      FollowersToJSON(user.Followers),
		Following:      FollowersToJSON(user.Following),
	}
}

func FollowersToJSON(followers []usermodel.UserFollowers) []UserFollowers {
	var convertedFollowers []UserFollowers
	for _, follower := range followers {
		convertedFollowers = append(convertedFollowers, UserFollowers{
			ID:       follower.ID,
			UserID:   follower.UserID,
			UserName: follower.UserName,
		})
	}
	return convertedFollowers
}
