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

// *******************************************************************
// WEIGHT
// *******************************************************************

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

	err := s.UserDAO.AddWeight(userID, weight, time.Now().Truncate(time.Hour*24), ctxLog)
	if err != nil {
		return err
	}

	return nil
}

// *******************************************************************
// BODY FAT
// *******************************************************************

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

func (s statService) AddBodyFat(userID string, bodyFat float32, ctxLog *log.Entry) error {

	ctxLog.Debugf("STATS_SERVICE: Processing AddBodyFat request for user: %s", userID)

	err := s.UserDAO.AddBodyFat(userID, bodyFat, time.Now().Truncate(time.Hour*24), ctxLog)
	if err != nil {
		return err
	}

	return nil
}

// *******************************************************************
// GYM ATTENDANCES (STREAK)
// *******************************************************************

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

func monday() time.Time {
	today := time.Now()
	todayIndex := int(today.Weekday())
	if todayIndex == 0 { // Sunday is 0 in Weekday(), make it equivalent to 7 for easier math
		todayIndex = 7
	}
	monday := today.AddDate(0, 0, -todayIndex+1).Truncate(24 * time.Hour).Add(-1) // Offset to Monday

	return monday
}

func (s statService) AddGymAttendance(userID string, date time.Time, ctxLog *log.Entry) error {

	ctxLog.Debugf("STATS_SERVICE: Processing AddGymAttendance request for user: %s", userID)

	// ===== Update current week =====
	monday := monday()

	// Calculate distantce to monday to know if date is in the current week
	if date.After(monday) {
		dateIndex := int(date.Sub(monday).Hours() / 24)
		if err := s.UserDAO.AddDayToCurrentWeek(userID, dateIndex, ctxLog); err != nil {
			return err
		}
	}

	// ===== Update gym attendances =====
	return s.UserDAO.AddGymAttendance(userID, date, ctxLog)
}

func (s statService) DeleteGymAttendance(userID string, date time.Time, ctxLog *log.Entry) error {

	ctxLog.Debugf("STATS_SERVICE: Processing DeleteGymAttendance request for user: %s", userID)

	// ===== Update current week =====
	monday := monday()

	// Calculate distantce to monday to know if date is in the current week
	if date.After(monday) {
		dateIndex := int(date.Sub(monday).Hours() / 24)
		if err := s.UserDAO.DeleteDayFromCurrentWeek(userID, dateIndex, ctxLog); err != nil {
			return err
		}
	}

	// ===== Update gym attendances =====
	return s.UserDAO.DeleteGymAttendance(userID, date, ctxLog)
}
