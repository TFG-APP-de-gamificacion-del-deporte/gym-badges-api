package badge_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IBadgeService interface {
	GetBadgesByUserID(userID string, ctxLog *log.Entry) (models.BadgesByUserResponse, error)
	AddBadge(userID string, badgeID int16, ctxLog *log.Entry) error
	DeleteBadge(userID string, badgeID int16, ctxLog *log.Entry) error
}
