package stats_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IStatsService interface {
	GetWeightHistory(userID string, months int32, ctxLog *log.Entry) (*models.MeasurementHistoryResponse, error)
	AddWeight(userID string, weight float32, ctxLog *log.Entry) error
	GetFatHistory(userID string, months int32, ctxLog *log.Entry) (*models.MeasurementHistoryResponse, error)
	AddBodyFat(userID string, fat float32, ctxLog *log.Entry) error
	GetStreakCalendarByYearAndMonth(userID string, year int32, month int32, ctxLog *log.Entry) (*models.StreakCalendarResponse, error)
}
