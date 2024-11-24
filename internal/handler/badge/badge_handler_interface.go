package badge_handler

import (
	"gym-badges-api/restapi/operations/badges"

	"github.com/go-openapi/runtime/middleware"
)

type IBadgeHandler interface {
	GetBadgesByUserID(params badges.GetBadgesByUserIDParams) middleware.Responder
}
