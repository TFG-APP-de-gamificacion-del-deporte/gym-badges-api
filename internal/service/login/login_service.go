package login_service

import (
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

func (s LoginService) Login(username, password string, ctxLog *log.Entry) (*models.LoginResponse, error) {

	ctxLog.Debugf("LOGIN_SERVICE: Processing login request for user: %s", username)

	user, err := s.userDAO.GetUser(username, ctxLog)
	if err != nil {
		return nil, err // FIXME Si el error es record not found, devolver BuildUnauthorizedError en vez de InternalServerError
	}

	if user.Password != password {
		return nil, customErrors.BuildUnauthorizedError("Invalid username or password")
	}

	token, err := s.sessionService.GenerateSession(username)
	if err != nil {
		return nil, err
	}

	response := models.LoginResponse{
		Token: token,
	}

	return &response, nil
}
