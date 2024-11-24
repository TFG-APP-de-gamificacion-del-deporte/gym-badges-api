package badge_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IBadgeService interface {
	GetBadgesByUserID(userID string, ctxLog *log.Entry) (*models.BadgesByUserResponse, error)
}
