package badge_dao

import (
	log "github.com/sirupsen/logrus"
)

type IBadgeDAO interface {
	GetBadges(ctxLog *log.Entry) ([]*Badge, error)
}
