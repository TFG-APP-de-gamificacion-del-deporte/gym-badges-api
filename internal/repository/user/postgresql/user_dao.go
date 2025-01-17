package postgresql

import (
	"database/sql"
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
		Preload("Preferences", func(db *gorm.DB) *gorm.DB {
			return db.Order("preference.id")
		}).
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

	if newUserInfo.Email != "" {
		user.Email = newUserInfo.Email
	}
	if newUserInfo.Name != "" {
		user.Name = newUserInfo.Name
	}
	if newUserInfo.Image != nil {
		user.Image = newUserInfo.Image
	}
	if newUserInfo.WeeklyGoal >= 1 && newUserInfo.WeeklyGoal <= 7 {
		user.WeeklyGoal = newUserInfo.WeeklyGoal
	}

	// Update top feats
	if newUserInfo.TopFeats != nil {
		if err := dao.connection.Model(&user).Association("TopFeats").Clear(); err != nil {
			return nil, err
		}
		if err := dao.connection.Model(&user).Association("TopFeats").Append(newUserInfo.TopFeats); err != nil {
			return nil, err
		}
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

func (dao userDAO) setDayToCurrentWeek(userID string, dayIndex int, marked bool, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Adding a day to current week of user: %s", userID)

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

	user.CurrentWeek[dayIndex] = marked

	if err := dao.connection.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func countDays(list []bool) int32 {
	var count int32 = 0
	for _, b := range list {
		if b {
			count++
		}
	}
	return count
}

func (dao userDAO) AddDayToCurrentWeek(userID string, dayIndex int, ctxLog *log.Entry) error {
	user, err := dao.setDayToCurrentWeek(userID, dayIndex, true, ctxLog)
	if err != nil {
		return err
	}

	// Update streak
	if countDays(user.CurrentWeek) == user.WeeklyGoal {
		user.Streak += 1
		if err := dao.connection.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func (dao userDAO) DeleteDayFromCurrentWeek(userID string, dayIndex int, ctxLog *log.Entry) error {
	user, err := dao.setDayToCurrentWeek(userID, dayIndex, false, ctxLog)
	if err != nil {
		return err
	}

	// Update streak
	if countDays(user.CurrentWeek) == user.WeeklyGoal-1 {
		user.Streak -= 1
		if err := dao.connection.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
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

	user.Weight = &weight
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

	user.BodyFat = &bodyFat
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
		Preload("GymAttendance").
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return queryResult.Error
	}

	user.GymAttendance = append(user.GymAttendance, userModelDB.GymAttendance{Date: date})

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

func (dao userDAO) GetAttendanceCount(userID string, ctxLog *log.Entry) (int32, error) {

	ctxLog.Debugf("USER_DAO: Getting attendance count for user: %s", userID)

	if err := dao.connection.Error; err != nil {
		return -1, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return -1, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return -1, queryResult.Error
	}

	count := dao.connection.Model(&user).Association("GymAttendance").Count()

	return int32(count), nil
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
		Preload("Preferences", func(db *gorm.DB) *gorm.DB {
			return db.Order("preference.id")
		}).
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

	const HIDE_STATS_ID = 1
	for _, f := range user.Friends {
		if f.Preferences[HIDE_STATS_ID].On {
			f.Weight = nil
			f.BodyFat = nil
		}
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

	queryResult = dao.connection.
		Where("id = ?", friendID).
		First(&friend)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, queryResult.Error
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

func (dao userDAO) GetFriendsCount(userID string, ctxLog *log.Entry) (int32, error) {

	ctxLog.Debugf("USER_DAO: Getting friends count for user: %s", userID)

	if err := dao.connection.Error; err != nil {
		return -1, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return -1, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return -1, queryResult.Error
	}

	var count int64

	queryResult = dao.connection.
		Model(&user).
		Joins(`JOIN user_friends ON "user".id = user_friends.friend_id OR "user".id = user_friends.user_id`).
		Where("user_friends.user_id = ? OR user_friends.friend_id = ?", userID, userID).
		Where(`"user".id != ?`, userID).
		Count(&count)

	if queryResult.Error != nil && errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return -1, queryResult.Error
	}

	return int32(count), nil
}

func (dao userDAO) CheckFriendship(userID string, friendID string, ctxLog *log.Entry) (bool, error) { // Returns: areFriends, confirmed, error

	ctxLog.Debugf("USER_DAO: Checking %s and %s friendship.", userID, friendID)

	if err := dao.connection.Error; err != nil {
		return false, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return false, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return false, queryResult.Error
	}

	var friend userModelDB.User

	queryResult = dao.connection.
		Where("id = ?", friendID).
		First(&friend)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return false, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return false, queryResult.Error
	}

	var count int64

	queryResult = dao.connection.
		Model(&user).
		Joins(`JOIN user_friends ON "user".id = user_friends.friend_id OR "user".id = user_friends.user_id`).
		Where("(user_friends.user_id = ? AND user_friends.friend_id = ?) OR (user_friends.friend_id = ? AND user_friends.user_id = ?)", userID, friendID, userID, friendID).
		Count(&count)

	if queryResult.Error != nil {
		return false, queryResult.Error
	}

	return count > 0, nil
}

func (dao userDAO) AddFriendRequest(userID string, friendID string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Adding friend request for %s (user) and %s (friend).", userID, friendID)

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

	queryResult = dao.connection.
		Where("id = ?", friendID).
		First(&friend)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, queryResult.Error
	}

	friend.FriendRequests = append(friend.FriendRequests, &user)

	return &friend, dao.connection.Save(&friend).Error
}

func (dao userDAO) GetUserWithFriendRequests(userID string, ctxLog *log.Entry) (*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting friend requests for user: %s ", userID)

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
		Preload("FriendRequests").
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return &user, nil
}

func (dao userDAO) CheckFriendRequest(userID string, friendID string, ctxLog *log.Entry) (bool, error) {

	ctxLog.Debugf("USER_DAO: Cheking if friendship request from %s exists for %s.", friendID, userID)

	if err := dao.connection.Error; err != nil {
		return false, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return false, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return false, queryResult.Error
	}

	var count int64

	queryResult = dao.connection.
		Model(user).
		Joins(`JOIN friend_requests ON "user".id = friend_requests.user_id`).
		Where("friend_requests.friend_request_id = ?", friendID).
		Count(&count)

	if queryResult.Error != nil && errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return false, queryResult.Error
	}

	return count > 0, nil
}

func (dao userDAO) DeleteFriendRequest(userID string, friendID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Adding friend request for %s (user) and %s (friend).", userID, friendID)

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

	if err := dao.connection.Unscoped().Model(&user).Association("FriendRequests").Unscoped().Delete(&friend); err != nil {
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

// *******************************************************************
// EXPERIENCE
// *******************************************************************

func (dao userDAO) AddExperience(userID string, exp int64, ctxLog *log.Entry) error {

	ctxLog.Debugf("USER_DAO: Adding experience to user: %s", userID)

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

	user.Experience += exp

	// Exp cannot get negative
	user.Experience = max(user.Experience, 0)

	return dao.connection.Save(&user).Error
}

// *******************************************************************
// RANKINGS
// *******************************************************************

func (dao *userDAO) GetUsersOrderedByExp(offset int64, size int32, ctxLog *log.Entry) ([]*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting global ranking offset: %d size: %d", offset, size)

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var users = make([]*userModelDB.User, 0)

	queryResult := dao.connection.
		Order("experience DESC, streak DESC, weekly_goal DESC, id").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&users)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return users, nil
}

func (dao *userDAO) GetUserWithGlobalRank(userID string, ctxLog *log.Entry) (*userModelDB.User, int64, error) {

	ctxLog.Debugf("USER_DAO: Getting user: %s with his global rank", userID)

	if err := dao.connection.Error; err != nil {
		return nil, -1, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, -1, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, -1, queryResult.Error
	}

	var rank int64 // 1-based
	queryResult = dao.connection.
		Raw(`
			SELECT rank 
			FROM (
				SELECT id, ROW_NUMBER() OVER (ORDER BY experience DESC) AS rank
				FROM "user"
			)
			WHERE "id" = ?
		`, userID).
		Scan(&rank)

	if queryResult.Error != nil {
		return nil, -1, queryResult.Error
	}

	return &user, rank, nil
}

func (dao *userDAO) GetFriendsOrderedByExp(userID string, offset int64, size int32, ctxLog *log.Entry) ([]*userModelDB.User, error) {

	ctxLog.Debugf("USER_DAO: Getting friends ranking for user: %s offset: %d size: %d", userID, offset, size)

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

	var friends = make([]*userModelDB.User, 0)

	queryResult = dao.connection.
		Order("experience DESC, streak DESC, weekly_goal DESC, id, name, image").
		Distinct("experience, streak, weekly_goal, id, name, image").
		Joins(`JOIN user_friends ON "user".id = user_friends.friend_id OR "user".id = user_friends.user_id`).
		Where("user_friends.user_id = ? OR user_friends.friend_id = ?", userID, userID).
		Limit(int(size)).
		Offset(int(offset)).
		Find(&friends)

	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return friends, nil
}

func (dao *userDAO) GetUserWithFriendsRank(userID string, ctxLog *log.Entry) (*userModelDB.User, int64, error) {

	ctxLog.Debugf("USER_DAO: Getting user: %s with his friends rank", userID)

	if err := dao.connection.Error; err != nil {
		return nil, -1, err
	}

	var user userModelDB.User

	queryResult := dao.connection.
		Where("id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, -1, customErrors.BuildNotFoundError(userNotFoundErrorMsg)
		}
		return nil, -1, queryResult.Error
	}

	var rank int64 // 1-based
	queryResult = dao.connection.
		Raw(`
			SELECT rank 
			FROM (
				SELECT id, ROW_NUMBER() OVER (ORDER BY experience DESC, streak DESC, weekly_goal DESC, id) AS rank
				FROM (
					SELECT DISTINCT ON ("user".id) "user".id, "user".experience, "user".streak, "user".weekly_goal
					FROM "user"
					JOIN user_friends 
						ON "user".id = user_friends.user_id
						OR "user".id = user_friends.friend_id 
					WHERE (user_friends.user_id = @user_id OR user_friends.friend_id = @user_id)
				)
			)
			WHERE "id" = @user_id
		`, sql.Named("user_id", userID)).
		Scan(&rank)

	if queryResult.Error != nil {
		return nil, -1, queryResult.Error
	}

	return &user, rank, nil
}
