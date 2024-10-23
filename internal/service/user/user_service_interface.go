package user_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	GetUser(userId string, ctxLog *log.Entry) (*models.GetUserInfoResponse, error)
}
