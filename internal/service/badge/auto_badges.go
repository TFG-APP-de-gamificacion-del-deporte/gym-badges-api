package badge_service

import (
	log "github.com/sirupsen/logrus"
)

var (
	streakBadges = []struct {
		badgeID int16
		weeks   int32
	}{
		{66, 1},  // BadgeID 66: The first week
		{67, 4},  // BadgeID 67: 4 weeks streak
		{68, 10}, // BadgeID 68: 10 weeks streak
		{69, 20}, // BadgeID 69: 20 weeks streak
		{70, 50}, // BadgeID 70: 50 weeks streak
	}
)

func (s badgesService) checkStreakBadges(userID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Checking streak badges.")

	user, err := s.userDAO.GetUser(userID, ctxLog)
	if err != nil {
		return err
	}

	for _, b := range streakBadges {
		if user.Streak >= b.weeks {
			hasBadge, err := s.badgeDAO.CheckBadge(userID, b.badgeID, ctxLog)
			if err != nil {
				return err
			}

			if !hasBadge {
				if err := s.badgeDAO.AddBadge(userID, b.badgeID, ctxLog); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
