package user

import (
	log "github.com/sirupsen/logrus"
)

type IUserDAO interface {
	GetUser(userID string, ctxLog *log.Entry) (*User, error)
	GetUserByEmail(email string, ctxLog *log.Entry) (*User, error)
	CreateUser(user *User, ctxLog *log.Entry) error
	GetUserWithWeightHistory(userID string, months int32, ctxLog *log.Entry) (*User, error)
	GetUserWithFatHistory(userID string, months int32, ctxLog *log.Entry) (*User, error)
	GetUserWithAttendance(userID string, year int32, month int32, ctxLog *log.Entry) (*User, error)
	GetUserWithFriends(userID string, offset int32, size int32, ctxLog *log.Entry) (*User, error)
	GetUserWithBadges(userID string, ctxLog *log.Entry) (*User, error)
	AddFriend(userID string, friendID string, ctxLog *log.Entry) (*User, error)
}
