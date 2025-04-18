package friends_handler

import (
	"gym-badges-api/restapi/operations/friends"

	"github.com/go-openapi/runtime/middleware"
)

type IFriendsHandler interface {
	GetFriendsByUserID(params friends.GetFriendsByUserIDParams) middleware.Responder
	AddFriend(params friends.AddFriendParams) middleware.Responder
	DeleteFriend(params friends.DeleteFriendParams) middleware.Responder
	GetFriendRequestsByUserID(params friends.GetFriendRequestsByUserIDParams) middleware.Responder
}
