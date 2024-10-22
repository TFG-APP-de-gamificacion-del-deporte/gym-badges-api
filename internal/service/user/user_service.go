package user_service

import (
	userDAO "gym-badges-api/internal/repository/user"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

func NewUserService(userDAO userDAO.IUserDAO) IUserService {
	return &UserService{
		UserDAO: userDAO,
	}
}

type UserService struct {
	UserDAO userDAO.IUserDAO
}

func (UserService) GetUser(userId string, ctxLog *log.Entry) (*models.GetUserInfoResponse, error) {
	return nil, nil
}
