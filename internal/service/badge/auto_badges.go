package badge_service

import (
	"time"

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

	attendancesBadges = []struct {
		badgeID     int16
		attendances int32
	}{
		{72, 100}, // BadgeID 72: Reach 100 gym sessions
		{73, 200}, // BadgeID 73: Reach 200 gym sessions
		{74, 500}, // BadgeID 74: Reach 500 gym sessions
	}

	timeBadges = []struct {
		badgeID int16
		time    time.Duration
	}{
		{71, time.Hour * 24 * 30},      // BadgeID 71: First month
		{75, time.Hour * 24 * 182},     // BadgeID 75: Six months
		{76, time.Hour * 24 * 365},     // BadgeID 76: First year!
		{77, time.Hour * 24 * 365 * 2}, // BadgeID 77: Two years
		{78, time.Hour * 24 * 365 * 3}, // BadgeID 78: Three years
		{79, time.Hour * 24 * 365 * 5}, // BadgeID 79: Five years
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

func (s badgesService) checkAttendancesBadges(userID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Checking attendances badges.")

	attendanceCount, err := s.userDAO.GetAttendanceCount(userID, ctxLog)
	if err != nil {
		return err
	}

	for _, b := range attendancesBadges {
		if attendanceCount >= b.attendances {
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

func (s badgesService) checkTimeBadges(userID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Checking time badges.")

	user, err := s.userDAO.GetUser(userID, ctxLog)
	if err != nil {
		return err
	}

	userLifespan := time.Since(user.CreatedAt)

	for _, b := range timeBadges {
		if userLifespan >= b.time {
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
