package stats_service

import (
	"gym-badges-api/internal/constants"
	userDAO "gym-badges-api/internal/repository/user"
	sessionService "gym-badges-api/internal/service/session"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

func NewStatsService(userDAO userDAO.IUserDAO, sessionService sessionService.ISessionService) IStatsService {
	return &statService{
		UserDAO:        userDAO,
		sessionService: sessionService,
	}
}

type statService struct {
	UserDAO        userDAO.IUserDAO
	sessionService sessionService.ISessionService
}

func (s statService) GetWeightHistory(userID string, months int32, ctxLog *log.Entry) (*models.GetWeightHistoryResponse, error) {

	ctxLog.Debugf("STATS_SERVICE: Processing GetWeightHistory request for user: %s", userID)

	user, err := s.UserDAO.GetUserWithWeightHistory(userID, months, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.GetWeightHistoryResponse{
		Days: make([]*models.WeightPerDay, 0, len(user.WeightHistory)),
	}

	for _, weight := range user.WeightHistory {
		response.Days = append(response.Days, &models.WeightPerDay{
			Date:  weight.Date.Format(constants.ISODateLayout),
			Value: weight.Weight,
		})
	}

	return &response, nil
}
