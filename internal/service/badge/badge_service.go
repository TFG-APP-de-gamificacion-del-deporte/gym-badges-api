package badge_service

import (
	"fmt"
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func NewBadgeService(userDAO userDAO.IUserDAO, badgeDAO badgeDAO.IBadgeDAO) IBadgeService {
	return &badgesService{
		userDAO:  userDAO,
		badgeDAO: badgeDAO,
	}
}

type badgesService struct {
	userDAO  userDAO.IUserDAO
	badgeDAO badgeDAO.IBadgeDAO
}

func (s badgesService) GetBadgesByUserID(userID string, ctxLog *log.Entry) (*models.BadgesByUserResponse, error) {

	ctxLog.Debugf("BADGES_SERVICE: Processing GetBadgesByUserID for user: %s", userID)

	var (
		eg     *errgroup.Group
		user   *userDAO.User
		badges []*badgeDAO.Badge
	)

	eg = new(errgroup.Group)

	eg.Go(func() error {
		var err error
		user, err = s.userDAO.GetUserWithBadges(userID, ctxLog)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		var err error
		badges, err = s.badgeDAO.GetBadges(ctxLog)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	userBadgesMap := make(map[uint16]bool)
	for _, badge := range user.Badges {
		userBadgesMap[badge.ID] = true
	}

	response := models.BadgesByUserResponse{
		Badges: make(map[string]models.Badge),
	}

	for _, badge := range badges {

		response.Badges[fmt.Sprint(badge.ID)] = models.Badge{
			Achieved:    userBadgesMap[badge.ID],
			Description: badge.Description,
			Image:       badge.Image,
			Name:        badge.Name,
			Parent:      fmt.Sprint(badge.ParentBadgeID),
		}
	}

	return &response, nil
}
