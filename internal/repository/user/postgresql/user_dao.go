package postgresql

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/internal/repository/config/postgresql"
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
	connection := postgresql.OpenConnection()
	return &userDAO{connection: connection}
}

// *******************************************************************
// USER
// *******************************************************************

func (dao userDAO) GetUser(userID string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting user: %s", userID)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Preload("TopFeats", func(db *gorm.DB) *gorm.DB {
			return db.Limit(3)
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

func (dao userDAO) EditUserInfo(userID string, newUserInfo *userModelDB.User, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Editing information of user: %s", userID)

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

	user.Email = newUserInfo.Email
	user.Name = newUserInfo.Name
	user.Image = newUserInfo.Image
	user.WeeklyGoal = newUserInfo.WeeklyGoal

	// Update top feats
	if err := dao.connection.Model(&user).Association("TopFeats").Clear(); err != nil {
		return nil, err
	}
	if err := dao.connection.Model(&user).Association("TopFeats").Append(newUserInfo.TopFeats); err != nil {
		return nil, err
	}

	// Update preferences
	if err := dao.connection.Unscoped().Model(&user).Association("Preferences").Unscoped().Delete(newUserInfo.Preferences); err != nil {
		return nil, err
	}
	if err := dao.connection.Model(&user).Association("Preferences").Append(newUserInfo.Preferences); err != nil {
		return nil, err
	}

	if err := dao.connection.Save(&user).Error; err != nil {
		return nil, err
	}

	user.Preferences = newUserInfo.Preferences

	return &user, nil
}

// *******************************************************************
// WEIGHT
// *******************************************************************

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

func (dao userDAO) AddWeight(userID string, weight float32, date time.Time, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Adding new weight to user %s", userID)

	if err := dao.connection.Error; err != nil {
		return err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	user.Weight = weight
	err := dao.connection.Unscoped().Model(&user).Association("WeightHistory").Unscoped().Delete(&userModelDB.WeightHistory{UserID: user.ID, Date: date})
	if err != nil {
		return err
	}
	err = dao.connection.Model(&user).Association("WeightHistory").Append(&userModelDB.WeightHistory{Date: date, Weight: weight})
	if err != nil {
		return err
	}

	return dao.connection.Save(&user).Error
}

// *******************************************************************
// BODY FAT
// *******************************************************************

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

func (dao userDAO) AddBodyFat(userID string, bodyFat float32, date time.Time, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Adding new body fat to user %s", userID)

	if err := dao.connection.Error; err != nil {
		return err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	user.BodyFat = bodyFat
	err := dao.connection.Unscoped().Model(&user).Association("FatHistory").Unscoped().Delete(&userModelDB.FatHistory{UserID: user.ID, Date: date})
	if err != nil {
		return err
	}
	err = dao.connection.Model(&user).Association("FatHistory").Append(&userModelDB.FatHistory{Date: date, Fat: bodyFat})
	if err != nil {
		return err
	}

	return dao.connection.Save(&user).Error
}

// *******************************************************************
// GYM ATTENDANCES (STREAK)
// *******************************************************************

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

func (dao userDAO) AddGymAttendance(userID string, date time.Time, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Adding a gym attendance to user %s", userID)

	if err := dao.connection.Error; err != nil {
		return err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	err := dao.connection.Model(&user).Association("GymAttendance").Append(&userModelDB.GymAttendance{Date: date})
	if err != nil {
		return err
	}

	return dao.connection.Save(&user).Error
}

func (dao userDAO) DeleteGymAttendance(userID string, date time.Time, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Deleting a gym attendance to user %s", userID)

	if err := dao.connection.Error; err != nil {
		return err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	err := dao.connection.Unscoped().Model(&user).Association("GymAttendance").Unscoped().Delete(&userModelDB.GymAttendance{
		UserID: user.ID,
		Date:   date,
	})
	if err != nil {
		return err
	}

	return dao.connection.Save(&user).Error
}

// *******************************************************************
// FRIENDS
// *******************************************************************

func (dao userDAO) GetUserWithFriends(userID string, offset int32, size int32, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting friends for user: %s offset: %d size: %d", userID, offset, size)

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

	queryResult = dao.connection.
		Preload("TopFeats", func(db *gorm.DB) *gorm.DB {
			return db.Limit(3)
		}).
		Joins(`JOIN user_friends ON "user".id = user_friends.friend_id OR "user".id = user_friends.user_id`).
		Where("user_friends.user_id = ? OR user_friends.friend_id = ?", userID, userID).
		Where(`"user".id != ?`, userID).
		Limit(int(size)).
		Offset(int(offset)).
		Find(&user.Friends)

	if queryResult.Error != nil && errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return nil, queryResult.Error
	}

	return &user, nil
}

func (dao userDAO) AddFriend(userID string, friendID string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Making %s (user) and %s (friend) friends.", userID, friendID)

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

	var friend userModelDB.User

	friendQueryResult := dao.connection.
		Where("id = ?", friendID).
		First(&friend)

	if friendQueryResult.Error != nil {
		if errors.Is(friendQueryResult.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, friendQueryResult.Error
	}

	user.Friends = append(user.Friends, &friend)
	if err := dao.connection.Save(user).Error; err != nil {
		return nil, err
	}

	return &friend, nil
}

func (dao userDAO) DeleteFriend(userID string, friendID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Making %s (user) and %s (friend) no longer friends.", userID, friendID)

	if err := dao.connection.Error; err != nil {
		return err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	var friend userModelDB.User

	queryResult = dao.connection.
		Where("id = ?", friendID).
		First(&friend)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	// Try to delete from both users friends list (The friendship relation will only be in one of them)
	if err := dao.connection.Model(&user).Association("Friends").Delete(&friend); err != nil {
		return err
	}
	if err := dao.connection.Model(&friend).Association("Friends").Delete(&user); err != nil {
		return err
	}

	return nil
}

// *******************************************************************
// BADGES
// *******************************************************************

func (dao userDAO) GetUserWithBadges(userID string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting badges for user: %s", userID)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Preload("Badges").
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
