package login_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type ILoginService interface {
	Login(username, password string, ctxLog *log.Entry) (*models.LoginResponse, error)
}
