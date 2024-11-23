package stats_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IStatsService interface {
	GetWeightHistory(userID string, months int32, ctxLog *log.Entry) (*models.MeasurementHistoryResponse, error)
}
