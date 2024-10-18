package login_service

import (
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	"gym-badges-api/models"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func NewLoginService(userDAO userDAO.IUserDAO) ILoginService {
	return &LoginService{
		UserDAO: userDAO,
	}
}

type LoginService struct {
	UserDAO userDAO.IUserDAO
}

func (s LoginService) Login(username, password string, ctxLog *log.Entry) (*models.LoginResponse, error) {

	ctxLog.Debug("LOGIN_SERVICE: Processing login request for user: %s", username)

	user, err := s.UserDAO.GetUser(username, ctxLog)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, customErrors.BuildUnauthorizedError("Invalid username or password")
	}

	response := models.LoginResponse{
		Token: uuid.New().String(),
	}

	return &response, nil
}
