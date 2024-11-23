package friends

import (
	"gym-badges-api/restapi/operations/friends"

	"github.com/go-openapi/runtime/middleware"
)

type IFriendsHandler interface {
	GetFriendsByUserID(params friends.GetFriendsByUserIDParams) middleware.Responder
}
