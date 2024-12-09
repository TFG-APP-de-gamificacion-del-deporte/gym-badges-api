package user_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	sessionService "gym-badges-api/internal/service/session"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
		UserID:      user.ID,
		BodyFat:     user.BodyFat,
		CurrentWeek: user.CurrentWeek,
		Experience:  user.Experience,
		Image:       user.Image,
		Name:        user.Name,
		Streak:      user.Streak,
		Weight:      user.Weight,
		TopFeats:    mapTopFeats(user.TopFeats),
	}

	return &response, nil
}

func mapTopFeats(dbTopFeats []*badgeDAO.Badge) []*models.Feat {

	topFeats := make([]*models.Feat, len(dbTopFeats))

	for i, badge := range dbTopFeats {
		topFeats[i] = &models.Feat{
			ID:          int32(badge.ID),
			Description: badge.Description,
			Image:       badge.Image,
			Name:        badge.Name,
		}
	}

	return topFeats
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

	// Encrypt password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, err
	}
	hash := string(bytes)

	newUser := userDAO.User{
		ID:          user.UserID,
		BodyFat:     0,
		CurrentWeek: []bool{false, false, false, false, false, false, false},
		Email:       user.Email,
		Experience:  0,
		Image:       user.Image,
		Name:        user.Name,
		Password:    hash,
		Streak:      0,
		Weight:      0,
	}

	if err = s.UserDAO.CreateUser(&newUser, ctxLog); err != nil {
		return nil, err
	}

	token, err := s.sessionService.GenerateSession(newUser.ID)
	if err != nil {
		return nil, err
	}

	response := models.LoginResponse{
		Token: token,
	}

	return &response, nil
}
