package login_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	sessionService "gym-badges-api/internal/service/session"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

func NewLoginService(userDAO userDAO.IUserDAO, sessionService sessionService.ISessionService) ILoginService {
	return &LoginService{
		userDAO:        userDAO,
		sessionService: sessionService,
	}
}

type LoginService struct {
	userDAO        userDAO.IUserDAO
	sessionService sessionService.ISessionService
}

func (s LoginService) Login(userID, password string, ctxLog *log.Entry) (*models.LoginResponse, error) {

	ctxLog.Debugf("LOGIN_SERVICE: Processing login request for user: %s", userID)

	user, err := s.userDAO.GetUser(userID, ctxLog)
	if err != nil {
		if errors.As(err, &customErrors.NotFoundError{}) {
			return nil, customErrors.BuildUnauthorizedError("Invalid username or password")
		}
		return nil, err
	}

	if user.Password != password {
		return nil, customErrors.BuildUnauthorizedError("Invalid username or password")
	}

	token, err := s.sessionService.GenerateSession(userID)
	if err != nil {
		return nil, err
	}

	response := models.LoginResponse{
		Token: token,
	}

	return &response, nil
}
