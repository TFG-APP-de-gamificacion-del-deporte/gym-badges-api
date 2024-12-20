package postgresql

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	badgeModelDB "gym-badges-api/internal/repository/badge"
	"gym-badges-api/internal/repository/config/postgresql"
	userModelDB "gym-badges-api/internal/repository/user"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	userNotFoundErrorMsg  = "User not found"
	badgeNotFoundErrorMsg = "Badge not found"
)

type badgeDAO struct {
	connection *gorm.DB
}

func NewBadgeDAO() badgeModelDB.IBadgeDAO {
	connection := postgresql.OpenConnection()
	return &badgeDAO{connection: connection}
}

func (dao badgeDAO) GetBadges(ctxLog *log.Entry) ([]*badgeModelDB.Badge, error) {

	ctxLog.Debugf("BADGE_DAO: Getting available badges")

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var badges = make([]*badgeModelDB.Badge, 0)

	queryResult := dao.connection.
		Find(&badges)

	if queryResult.Error != nil && !errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return nil, queryResult.Error
	}

	return badges, nil
}

func (dao badgeDAO) GetBadge(badgeID int16, ctxLog *log.Entry) (*badgeModelDB.Badge, error) {

	ctxLog.Debugf("BADGE_DAO: Getting available badges")

	if err := dao.connection.Error; err != nil {
		return nil, err
	}

	var badge badgeModelDB.Badge

	queryResult := dao.connection.
		Where("id = ?", badgeID).
		First(&badge)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.BuildNotFoundError(badgeNotFoundErrorMsg)
		}
		return nil, queryResult.Error
	}

	return &badge, nil
}

func (dao badgeDAO) AddBadge(userID string, badgeID int16, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGE_DAO: Adding badge %d to user %s", badgeID, userID)

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

	user.Badges = append(user.Badges, &badgeModelDB.Badge{ID: badgeID})

	return dao.connection.Save(user).Error
}

func (dao badgeDAO) DeleteBadge(userID string, badgeID int16, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGE_DAO: Adding badge %d to user %s", badgeID, userID)

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

	return dao.connection.Unscoped().Model(&user).Association("Badges").Unscoped().Delete(badgeModelDB.Badge{ID: badgeID})
}

func (dao badgeDAO) CheckBadge(userID string, badgeID int16, ctxLog *log.Entry) (bool, error) {

	ctxLog.Debugf("BADGE_DAO: Checking user %s has badge %d", userID, badgeID)

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

	var badges []badgeModelDB.Badge
	dao.connection.Model(&user).Where("badge_id == ?", badgeID).Association("Badges").Find(&badges)

	if queryResult.Error != nil {
		return false, queryResult.Error
	}

	return (len(badges) > 0), nil
}
