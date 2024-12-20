package badge_service

import (
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (s badgesService) checkAutoBadges(userID string, ctxLog *log.Entry) error {

	eg := new(errgroup.Group)

	eg.Go(func() error {
		return s.checkStreakBadges(userID, ctxLog)
	})

	eg.Go(func() error {
		return s.checkAttendancesBadges(userID, ctxLog)
	})

	eg.Go(func() error {
		return s.checkTimeBadges(userID, ctxLog)
	})

	eg.Go(func() error {
		return s.checkGlobalRankingBadges(userID, ctxLog)
	})

	eg.Go(func() error {
		return s.checkFriendsRankingBadges(userID, ctxLog)
	})

	eg.Go(func() error {
		return s.checkFriendCountBadges(userID, ctxLog)
	})

	return eg.Wait()
}

// *******************************************************************
// CONSISTENCY BADGES
// *******************************************************************

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

// *******************************************************************
// RANKING BADGES
// *******************************************************************

var (
	globalRankingBadges = []struct {
		badgeID int16
		rank    int64
	}{
		{80, 500}, // BadgeID 80: Get to the top 500 at the Global Ranking
		{81, 100}, // BadgeID 81: Get to the top 100 at the Global Ranking
		{82, 50},  // BadgeID 82: Get to the top 50 at the Global Ranking
		{83, 10},  // BadgeID 83: Get to the top 10 at the Global Ranking
		{84, 3},   // BadgeID 84: Get to the top 3 at the Global Ranking
		{85, 1},   // BadgeID 85: Get to the top 1 at the Global Ranking
	}

	friendsRankingBadges = []struct {
		badgeID    int16
		minFriends int32
		rank       int64
	}{
		{90, 20, 3}, // BadgeID 90: Reach the podium at the Friends Ranking with at least 20 friends
		{91, 20, 1}, // BadgeID 91: Get to the top at the Friends Ranking with at least 20 friends
		{92, 10, 3}, // BadgeID 92: Reach the podium at the Friends Ranking with at least 10 friends
		{93, 10, 1}, // BadgeID 93: Get to the top at the Friends Ranking with at least 10 friends
		{94, 5, 1},  // BadgeID 94: Get to the top at the Friends Ranking with at least 5 friends
	}
)

func (s badgesService) checkGlobalRankingBadges(userID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Checking global ranking badges.")

	_, rank, err := s.userDAO.GetUserWithGlobalRank(userID, ctxLog)
	if err != nil {
		return err
	}

	for _, b := range globalRankingBadges {
		if rank <= b.rank {
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

func (s badgesService) checkFriendsRankingBadges(userID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Checking friends ranking badges.")

	_, rank, err := s.userDAO.GetUserWithFriendsRank(userID, ctxLog)
	if err != nil {
		return err
	}

	friendsCount, err := s.userDAO.GetFriendsCount(userID, ctxLog)
	if err != nil {
		return err
	}

	for _, b := range friendsRankingBadges {
		if friendsCount >= b.minFriends && rank <= b.rank {
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

// *******************************************************************
// FRIEND COUNT BADGES
// *******************************************************************

var (
	friendCountBadges = []struct {
		badgeID    int16
		minFriends int32
	}{
		{86, 1},  // BadgeID 86: Add your first friend
		{87, 5},  // BadgeID 87: Add five friends
		{88, 10}, // BadgeID 88: Add ten friends
		{89, 20}, // BadgeID 89: Add twenty friends
	}
)

func (s badgesService) checkFriendCountBadges(userID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Checking friend count badges.")

	friendsCount, err := s.userDAO.GetFriendsCount(userID, ctxLog)
	if err != nil {
		return err
	}

	for _, b := range friendCountBadges {
		if friendsCount >= b.minFriends {
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
