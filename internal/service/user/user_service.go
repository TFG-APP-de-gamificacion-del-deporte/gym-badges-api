package user_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	sessionService "gym-badges-api/internal/service/session"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

func NewUserService(userDAO userDAO.IUserDAO, sessionService sessionService.ISessionService) IUserService {
	return &UserService{
		UserDAO:        userDAO,
		sessionService: sessionService,
	}
}

type UserService struct {
	UserDAO        userDAO.IUserDAO
	sessionService sessionService.ISessionService
}

func (s UserService) GetUser(userID string, ctxLog *log.Entry) (*models.GetUserInfoResponse, error) {

	ctxLog.Debugf("USER_SERVICE: Processing getUserInfo request for user: %s", userID)

	user, err := s.UserDAO.GetUser(userID, ctxLog)

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

func (s UserService) CreateUser(user *models.CreateUserRequest, ctxLog *log.Entry) (*models.LoginResponse, error) {

	ctxLog.Debugf("USER_SERVICE: Processing user creation request: %s", user.UserID)

	userInDB, err := s.UserDAO.GetUser(user.UserID, ctxLog)
	if err != nil && !errors.As(err, &customErrors.NotFoundError{}) {
		return nil, err
	}

	if userInDB != nil {
		return nil, customErrors.BuildConflictError("user %s already exists", user.UserID)
	}

	userInDB, err = s.UserDAO.GetUserByEmail(user.Email, ctxLog)
	if err != nil && !errors.As(err, &customErrors.NotFoundError{}) {
		return nil, err
	}

	if userInDB != nil {
		return nil, customErrors.BuildConflictError("email %s already exists", user.Email)
	}

	newUser := userDAO.User{
		UserID:   user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}

	if err = s.UserDAO.CreateUser(&newUser, ctxLog); err != nil {
		return nil, err
	}

	token, err := s.sessionService.GenerateSession(newUser.UserID)
	if err != nil {
		return nil, err
	}

	response := models.LoginResponse{
		Token: token,
	}

	return &response, nil
}
