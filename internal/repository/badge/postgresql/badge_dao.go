package postgresql

import (
	"errors"
	badgeModelDB "gym-badges-api/internal/repository/badge"
	"gym-badges-api/internal/repository/config/postgresql"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
