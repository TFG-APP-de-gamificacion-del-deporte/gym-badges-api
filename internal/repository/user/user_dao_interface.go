package user

import (
	log "github.com/sirupsen/logrus"
)

type IUserDAO interface {
	GetUser(userID string, ctxLog *log.Entry) (*User, error)
}
