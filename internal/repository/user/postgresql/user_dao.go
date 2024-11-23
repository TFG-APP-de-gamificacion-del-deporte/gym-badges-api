package postgresql

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userModelDB "gym-badges-api/internal/repository/user"
	"time"

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

func (dao userDAO) GetUserWithWeightHistory(userID string, months int32, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting weight history for user: %s for last %d months", userID, months)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Preload("WeightHistory", func(db *gorm.DB) *gorm.DB {
			if months > 0 {
				startDate := time.Now().AddDate(0, -int(months), 0)
				return db.Where("date >= ?", startDate).Order("date ASC")
			}
			return db.Order("date ASC")
		}).
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

func (dao userDAO) GetUserWithFatHistory(userID string, months int32, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting fat history for user: %s for last %d months", userID, months)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Preload("FatHistory", func(db *gorm.DB) *gorm.DB {
			if months > 0 {
				startDate := time.Now().AddDate(0, -int(months), 0)
				return db.Where("date >= ?", startDate).Order("date ASC")
			}
			return db.Order("date ASC")
		}).
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

func (dao userDAO) GetUserWithAttendance(userID string, year int32, month int32, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting attendance info for user: %s in year %d and month %d", userID, year, month)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Preload("GymAttendance", func(db *gorm.DB) *gorm.DB {
			return db.Where("EXTRACT(YEAR FROM date) = ? AND EXTRACT(MONTH FROM date) = ?", year, month).
				Order("date ASC")
		}).
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
