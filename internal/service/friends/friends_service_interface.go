package friends_service

import (
	"gym-badges-api/models"

	log "github.com/sirupsen/logrus"
)

type IFriendsService interface {
	GetFriendsByUserID(userID string, page int32, ctxLog *log.Entry) (*models.FriendsResponse, error)
	AddFriend(userID string, friendID string, ctxLog *log.Entry) (*models.FriendInfo, error)
	DeleteFriend(userID string, friendID string, ctxLog *log.Entry) error
	GetFriendRequestsByUserID(userID string, ctxLog *log.Entry) (*models.FriendRequestsResponse, error)
}
