package friends_service

import (
	configs "gym-badges-api/config/gym-badges-server"
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	"gym-badges-api/models"
	"gym-badges-api/tools/utils"

	log "github.com/sirupsen/logrus"
)

func NewFriendsService(userDAO userDAO.IUserDAO) IFriendsService {
	return &friendsService{
		UserDAO: userDAO,
	}
}

type friendsService struct {
	UserDAO userDAO.IUserDAO
}

func (s friendsService) GetFriendsByUserID(userID string, page int32, ctxLog *log.Entry) (*models.FriendsResponse, error) {

	ctxLog.Debugf("FRIENDS_SERVICE: Processing GetFriendsByUserID for user: %s", userID)

	offset := (page - 1) * configs.Basic.FriendsPageSize
	size := configs.Basic.FriendsPageSize

	user, err := s.UserDAO.GetUserWithFriends(userID, offset, size, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.FriendsResponse{
		Friends: make([]*models.FriendInfo, len(user.Friends)),
	}

	for i, friend := range user.Friends {

		response.Friends[i] = &models.FriendInfo{
			Fat:      friend.BodyFat,
			Image:    friend.Image,
			Level:    utils.CalcLevel(friend.Experience),
			Name:     friend.Name,
			Streak:   friend.Streak,
			TopFeats: mapTopFeats(friend.TopFeats),
			User:     friend.ID,
			Weight:   friend.Weight,
		}
	}

	return &response, nil
}

func mapTopFeats(badges []*badgeDAO.Badge) []*models.Feat {

	topFeats := make([]*models.Feat, len(badges))

	for i, badge := range badges {
		topFeats[i] = &models.Feat{
			ID:          int32(badge.ID),
			Description: badge.Description,
			Image:       badge.Image,
			Name:        badge.Name,
		}
	}

	return topFeats
}

func (s friendsService) AddFriend(userID string, friendID string, ctxLog *log.Entry) (*models.FriendInfo, error) {

	ctxLog.Debugf("FRIENDS_SERVICE: Making %s (user) and %s (friend) friends.", userID, friendID)

	areFriends, err := s.UserDAO.CheckFriendship(userID, friendID, ctxLog)
	if err != nil {
		return nil, err
	}

	var friend *userDAO.User

	if !areFriends {
		friendshipRequested, err := s.UserDAO.CheckFriendRequest(userID, friendID, ctxLog)
		if err != nil {
			return nil, err
		}

		if !friendshipRequested {
			// Add friend request
			friend, err = s.UserDAO.AddFriendRequest(userID, friendID, ctxLog)
			if err != nil {
				return nil, err
			}
		} else {
			// Add friend
			friend, err = s.UserDAO.AddFriend(userID, friendID, ctxLog)
			if err != nil {
				return nil, err
			}
			// Delete friend request
			err = s.UserDAO.DeleteFriendRequest(userID, friendID, ctxLog)
			if err != nil {
				return nil, err
			}
		}
	} else {
		friend, err = s.UserDAO.GetUser(friendID, ctxLog)
		if err != nil {
			return nil, err
		}
	}

	friendInfo := models.FriendInfo{
		Fat:      friend.BodyFat,
		Image:    friend.Image,
		Level:    utils.CalcLevel(friend.Experience),
		Name:     friend.Name,
		Streak:   friend.Streak,
		TopFeats: mapTopFeats(friend.Badges),
		User:     friend.ID,
		Weight:   friend.Weight,
	}

	return &friendInfo, nil
}

func (s friendsService) DeleteFriend(userID string, friendID string, ctxLog *log.Entry) error {

	ctxLog.Debugf("FRIENDS_SERVICE: Making %s (user) and %s (friend) no longer friends.", userID, friendID)

	areFriends, err := s.UserDAO.CheckFriendship(userID, friendID, ctxLog)
	if err != nil {
		return err
	}

	if areFriends {
		return s.UserDAO.DeleteFriend(userID, friendID, ctxLog)
	} else {
		return s.UserDAO.DeleteFriendRequest(userID, friendID, ctxLog)
	}
}

func (s friendsService) GetFriendRequestsByUserID(userID string, ctxLog *log.Entry) (*models.FriendRequestsResponse, error) {

	ctxLog.Debugf("FRIENDS_SERVICE: Getting friend requests for user: %s ", userID)

	user, err := s.UserDAO.GetUserWithFriendRequests(userID, ctxLog)
	if err != nil {
		return nil, err
	}

	response := models.FriendRequestsResponse{
		FriendRequests: make([]*models.FriendRequestInfo, len(user.FriendRequests)),
	}

	for i, friendRequest := range user.FriendRequests {
		response.FriendRequests[i] = &models.FriendRequestInfo{
			User:  friendRequest.ID,
			Name:  friendRequest.Name,
			Image: friendRequest.Image,
		}
	}

	return &response, nil
}
