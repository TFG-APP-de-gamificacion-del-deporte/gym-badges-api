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

func (s UserService) GetUser(userId string, ctxLog *log.Entry) (*models.GetUserInfoResponse, error) {

	ctxLog.Debugf("USER_SERVICE: Processing getUserInfo request for user: %s", userId)

	user, err := s.UserDAO.GetUser(userId, ctxLog)

	if err != nil {
		return nil, err
	}

	response := models.GetUserInfoResponse{
		CurrentWeek: user.CurrentWeek,
		Experience:  user.Experience,
		BodyFat:     user.BodyFat,
		Image:       user.Image,
		Streak:      user.Streak,
		Weight:      user.Weight,
	}

	return &response, nil
}
