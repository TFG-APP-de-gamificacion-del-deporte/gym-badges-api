package postgresql

import (
	"errors"
	userModelDB "gym-badges-api/internal/repository/user"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		Where("user_id = ?", userID).
		First(&user)

	if queryResult.Error != nil {
		if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
			return nil, queryResult.Error
		}
		return nil, queryResult.Error
	}

	return &user, nil
}
