package user_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	GetUser(userID string, authUserID string, ctxLog *log.Entry) (*models.GetUserInfoResponse, error)
	CreateUser(request *models.CreateUserRequest, ctxLog *log.Entry) (*models.LoginResponse, error)
	EditUserInfo(userID string, request *models.EditUserInfoRequest, ctxLog *log.Entry) (*models.GetUserInfoResponse, error)
	EditUserPreferences(userID string, request *models.EditUserPreferenceRequest, ctxLog *log.Entry) error
	EditUserTopFeats(userID string, topFeats []int16, ctxLog *log.Entry) error
}
