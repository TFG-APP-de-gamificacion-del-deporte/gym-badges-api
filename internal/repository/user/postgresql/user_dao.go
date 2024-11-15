package postgresql

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userModelDB "gym-badges-api/internal/repository/user"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	userNotFoundErrorMsg = "User not found"
)

type userDAO struct {
	connection *gorm.DB
}

func NewUserDAO() userModelDB.IUserDAO {
	connection := OpenConnection()
	return &userDAO{connection: connection}
}

func (dao userDAO) GetUser(userID string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting user: %s", userID)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, queryResult.Error
	}

	return &user, nil
}

func (dao userDAO) GetUserByEmail(email string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting user by email: %s", email)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("email = ?", email).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, queryResult.Error
	}

	return &user, nil
}

func (dao userDAO) CreateUser(user *userModelDB.User, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Creating user: %s", user.ID)

	if err := dao.connection.Error; err != nil {
		return err
	}

	return dao.connection.Create(user).Error
}
