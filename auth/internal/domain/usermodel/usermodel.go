package usermodel

import (
	"time"
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
