package rankings_service

import (
	// configs "gym-badges-api/config/gym-badges-server"
	userDAO "gym-badges-api/internal/repository/user"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

func NewRankingsService(userDAO userDAO.IUserDAO) IRankingsService {
	return &rankingsService{
		// UserDAO: userDAO,
	}
}

type rankingsService struct {
	// UserDAO userDAO.IUserDAO
}

func (r *rankingsService) GetGlobalRanking(userID string, page int32, ctxLog *log.Entry) (*models.GetRankingResponse, error) {
	panic("unimplemented")
}
