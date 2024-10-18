package user

import (
	log "github.com/sirupsen/logrus"
)

type IUserDAO interface {
	GetUser(user string, ctxLog *log.Entry) (*User, error)
}
