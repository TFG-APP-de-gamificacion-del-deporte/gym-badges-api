package user

import (
	log "github.com/sirupsen/logrus"
)

type IUserDAO interface {
	GetUser(userID string, ctxLog *log.Entry) (*User, error)
	GetUserByEmail(email string, ctxLog *log.Entry) (*User, error)
	CreateUser(user *User, ctxLog *log.Entry) error
	GetUserWithWeightHistory(userID string, months int32, ctxLog *log.Entry) (*User, error)
}
