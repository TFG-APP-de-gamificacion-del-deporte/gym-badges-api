package badge_service

import (
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

func (s badgesService) GetBadgesByUserID(userID string, ctxLog *log.Entry) (models.BadgesByUserResponse, error) {

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

	response := make(models.BadgesByUserResponse, 0)

	badgeMap := make(map[int32][]*models.Badge)

	for _, badge := range badges {

		b := models.Badge{
			Achieved:    userBadgesMap[badge.ID],
			Description: badge.Description,
			ID:          int32(badge.ID),
			Image:       badge.Image,
			Name:        badge.Name,
		}

		if badge.ParentBadgeID == 0 {
			response = append(response, &b)
		} else {
			badgeMap[int32(badge.ParentBadgeID)] = append(badgeMap[int32(badge.ParentBadgeID)], &b)
		}
	}

	for i := range response {
		addChildren(response[i], badgeMap)
	}

	return response, nil
}

func addChildren(badge *models.Badge, auxMap map[int32][]*models.Badge) {

	children := auxMap[badge.ID]

	for i := range children {
		addChildren(children[i], auxMap)
	}

	badge.Children = children
}
