package badge_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	badgeService "gym-badges-api/internal/service/badge"
	"gym-badges-api/models"
	"gym-badges-api/restapi/operations/badges"
	op "gym-badges-api/restapi/operations/badges"
	toolsLogging "gym-badges-api/tools/logging"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

var (
	unauthorizedErrorResponse = models.GenericResponse{
		Code:    fmt.Sprint(http.StatusUnauthorized),
		Message: http.StatusText(http.StatusUnauthorized),
	}

	notFoundErrorResponse = models.GenericResponse{
		Code:    fmt.Sprint(http.StatusNotFound),
		Message: http.StatusText(http.StatusNotFound),
	}

	internalServerErrorResponse = models.GenericResponse{
		Code:    fmt.Sprint(http.StatusInternalServerError),
		Message: http.StatusText(http.StatusInternalServerError),
	}
)

func NewBadgeHandler(badgeService badgeService.IBadgeService) IBadgeHandler {
	return &badgesHandler{
		badgeService: badgeService,
	}
}

type badgesHandler struct {
	badgeService badgeService.IBadgeService
}

func (h badgesHandler) GetBadgesByUserID(params badges.GetBadgesByUserIDParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("BADGES_HANDLER: Getting badges for user: %s", params.UserID)

	response, err := h.badgeService.GetBadgesByUserID(params.UserID, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetBadgesByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		default:
			return op.NewGetBadgesByUserIDInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetBadgesByUserIDOK().WithPayload(response)
}
