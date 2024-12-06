package stats_service

import (
	"gym-badges-api/internal/constants"
	userDAO "gym-badges-api/internal/repository/user"
	sessionService "gym-badges-api/internal/service/session"
	"gym-badges-api/models"
	"time"

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

func (s statService) GetWeightHistory(userID string, months int32, ctxLog *log.Entry) (*models.MeasurementHistoryResponse, error) {

	ctxLog.Debugf("STATS_SERVICE: Processing GetWeightHistory request for user: %s", userID)

	user, err := s.UserDAO.GetUserWithWeightHistory(userID, months, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.MeasurementHistoryResponse{
		Days: make([]*models.MeasurementPerDay, 0, len(user.WeightHistory)),
	}

	for _, weight := range user.WeightHistory {
		response.Days = append(response.Days, &models.MeasurementPerDay{
			Date:  weight.Date.Format(constants.ISODateLayout),
			Value: weight.Weight,
		})
	}

	return &response, nil
}

func (s statService) AddWeight(userID string, weight float32, ctxLog *log.Entry) error {

	ctxLog.Debugf("STATS_SERVICE: Processing AddWeight request for user: %s", userID)

	err := s.UserDAO.AddWeight(userID, weight, time.Now(), ctxLog)
	if err != nil {
		return err
	}

	return nil
}

func (s statService) GetFatHistory(userID string, months int32, ctxLog *log.Entry) (*models.MeasurementHistoryResponse, error) {

	ctxLog.Debugf("STATS_SERVICE: Processing GetFatHistory request for user: %s", userID)

	user, err := s.UserDAO.GetUserWithFatHistory(userID, months, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.MeasurementHistoryResponse{
		Days: make([]*models.MeasurementPerDay, 0, len(user.FatHistory)),
	}

	for _, weight := range user.FatHistory {
		response.Days = append(response.Days, &models.MeasurementPerDay{
			Date:  weight.Date.Format(constants.ISODateLayout),
			Value: weight.Fat,
		})
	}

	return &response, nil
}

func (s statService) GetStreakCalendarByYearAndMonth(userID string, year int32, month int32,
	ctxLog *log.Entry) (*models.StreakCalendarResponse, error) {

	ctxLog.Debugf("STATS_SERVICE: Processing GetStreakCalendar request for user: %s", userID)

	user, err := s.UserDAO.GetUserWithAttendance(userID, year, month, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.StreakCalendarResponse{
		Days:       make([]string, 0, len(user.GymAttendance)),
		Streak:     user.Streak,
		WeeklyGoal: user.WeeklyGoal,
	}

	for _, at := range user.GymAttendance {
		response.Days = append(response.Days, at.Date.Format(constants.ISODateLayout))
	}

	return &response, nil
}
