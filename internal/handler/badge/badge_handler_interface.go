package badge_handler

import (
	"gym-badges-api/restapi/operations/badges"

	"github.com/go-openapi/runtime/middleware"
)

type IBadgeHandler interface {
	GetBadgesByUserID(params badges.GetBadgesByUserIDParams) middleware.Responder
	AddBadge(params badges.AddBadgeParams) middleware.Responder
	DeleteBadge(params badges.DeleteBadgeParams) middleware.Responder
}
