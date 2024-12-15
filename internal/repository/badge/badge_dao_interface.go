package badge_dao

import (
	log "github.com/sirupsen/logrus"
)

type IBadgeDAO interface {
	GetBadges(ctxLog *log.Entry) ([]*Badge, error)
	AddBadge(userID string, badgeID int16, ctxLog *log.Entry) error
	GetBadge(badgeID int16, ctxLog *log.Entry) (*Badge, error)
	DeleteBadge(userID string, badgeID int16, ctxLog *log.Entry) error
}
