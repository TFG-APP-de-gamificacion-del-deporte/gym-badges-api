package rankings_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IRankingsService interface {
	GetGlobalRanking(userID string, page int32, ctxLog *log.Entry) (*models.GetRankingResponse, error)
	GetFriendsRanking(userID string, page int32, ctxLog *log.Entry) (*models.GetRankingResponse, error)
}
