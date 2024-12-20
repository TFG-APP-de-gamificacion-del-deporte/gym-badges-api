package badge_service

import (
	customErrors "gym-badges-api/internal/custom-errors"
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

	if err := s.checkStreakBadges(userID, ctxLog); err != nil {
		return nil, err
	}

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

	userBadgesMap := make(map[int16]bool)
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
			Exp:         badge.Exp,
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

	if children == nil {
		badge.Children = make([]*models.Badge, 0)
	} else {
		badge.Children = children
	}
}

func (s badgesService) AddBadge(userID string, badgeID int16, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Processing AddBadge for user: %s", userID)

	user, err := s.userDAO.GetUserWithBadges(userID, ctxLog)
	if err != nil {
		return err
	}

	badge, err := s.badgeDAO.GetBadge(badgeID, ctxLog)
	if err != nil {
		return err
	}

	// Check user already has badge's parent
	hasParent := false
	for _, b := range user.Badges {
		if b.ID == badge.ParentBadgeID {
			hasParent = true
			break
		}
	}

	if !hasParent {
		return customErrors.BuildForbiddenError("Parent Badge %d is needed first to mark badge %d as completed.", badge.ParentBadgeID, badge.ID)
	}

	err = s.badgeDAO.AddBadge(userID, badgeID, ctxLog)
	if err != nil {
		return err
	}

	// Add exp to user
	return s.userDAO.AddExperience(userID, badge.Exp, ctxLog)
}

func (s badgesService) DeleteBadge(userID string, badgeID int16, ctxLog *log.Entry) error {

	ctxLog.Debugf("BADGES_SERVICE: Processing DeleteBadge for user: %s", userID)

	user, err := s.userDAO.GetUserWithBadges(userID, ctxLog)
	if err != nil {
		return err
	}

	badge, err := s.badgeDAO.GetBadge(badgeID, ctxLog)
	if err != nil {
		return err
	}

	// Check user doesn't have children of the badge
	hasChildren := false
	for _, b := range user.Badges {
		if b.ParentBadgeID == badge.ID {
			hasChildren = true
			break
		}
	}

	if hasChildren {
		return customErrors.BuildForbiddenError("Cannot delete badge %d because user has children badges.", badge.ID)
	}

	err = s.badgeDAO.DeleteBadge(userID, badgeID, ctxLog)
	if err != nil {
		return err
	}

	// Add exp to user
	return s.userDAO.AddExperience(userID, -badge.Exp, ctxLog)
}
