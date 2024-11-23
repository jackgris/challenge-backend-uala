package userdb

import (
	"time"

	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
)

type User struct {
	ID             string
	UserName       string
	Email          string
	Password       string
	FollowerCount  int
	FollowingCount int
	Salt           string
	Token          string
	DateCreated    time.Time
	EncodedDate    string
	Followers      []UserFollowers
	Following      []UserFollowers
}

type UserFollowers struct {
	ID       string
	UserID   string
	UserName string
}

func UserToModel(user User) usermodel.User {
	return usermodel.User{
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
		Followers:      FollowersToModel(user.Followers),
		Following:      FollowersToModel(user.Following),
	}
}

func FollowersToModel(followers []UserFollowers) []usermodel.UserFollowers {
	var convertedFollowers []usermodel.UserFollowers
	for _, follower := range followers {
		convertedFollowers = append(convertedFollowers, usermodel.UserFollowers{
			ID:       follower.ID,
			UserID:   follower.UserID,
			UserName: follower.UserName,
		})
	}
	return convertedFollowers
}
