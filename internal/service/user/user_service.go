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
		Height:      *user.Height,
		Sex:         user.Sex,
		WeeklyGoal:  user.WeeklyGoal,
		TopFeats:    mapTopFeats(user.TopFeats),
		Preferences: mapPreferences(user.Preferences),
	}

	return &response, nil
}

func mapPreferences(dbPreferences []userDAO.Preference) []*models.Preference {
	preferences := make([]*models.Preference, len(dbPreferences))

	for i, p := range dbPreferences {
		preferences[i] = &models.Preference{
			PreferenceID: int32(p.ID),
			On:           p.On,
		}
	}

	return preferences
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
		BodyFat:     nil,
		CurrentWeek: []bool{false, false, false, false, false, false, false},
		Email:       user.Email,
		Experience:  0,
		Image:       user.Image,
		Name:        user.Name,
		Password:    hash,
		Streak:      0,
		Weight:      nil,
		Height:      &user.Height,
		Sex:         user.Sex,
		WeeklyGoal:  3,
		Preferences: []userDAO.Preference{
			{ID: 1, On: false, UserID: user.UserID}, // Private account
			{ID: 2, On: false, UserID: user.UserID}, // Hide weight, fat, height and sex
		},
		Badges: []*badgeDAO.Badge{ // Base category badges
			{ID: -1},
			{ID: -2},
			{ID: -3},
			{ID: -4},
			{ID: -5},
			{ID: -6},
			{ID: -7},
		},
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

func (s UserService) EditUserInfo(userID string, request *models.EditUserInfoRequest, ctxLog *log.Entry) (*models.GetUserInfoResponse, error) {

	ctxLog.Debugf("USER_SERVICE: Editing information of user: %s", userID)

	var newUserInfo userDAO.User

	newUserInfo.Email = request.Email
	newUserInfo.Name = request.Name
	newUserInfo.Image = request.Image
	newUserInfo.WeeklyGoal = request.WeeklyGoal
	newUserInfo.Height = &request.Height
	newUserInfo.Sex = request.Sex

	// Top feats
	if request.TopFeats != nil {
		// Limit length to 3
		n := min(len(request.TopFeats), 3)
		newTopFeats := make([]*badgeDAO.Badge, n)
		for i := 0; i < n; i++ {
			newTopFeats[i] = &badgeDAO.Badge{ID: int16(request.TopFeats[i])}
		}
		newUserInfo.TopFeats = newTopFeats
	}

	// Preferences
	for _, p := range request.Preferences {
		newUserInfo.Preferences = append(newUserInfo.Preferences, userDAO.Preference{UserID: userID, ID: uint(p.PreferenceID), On: p.On})
	}

	user, err := s.UserDAO.EditUserInfo(userID, &newUserInfo, ctxLog)
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
		Height:      *user.Height,
		Sex:         user.Sex,
		TopFeats:    mapTopFeats(user.TopFeats),
		Preferences: mapPreferences(user.Preferences),
	}

	return &response, nil
}
