package badge_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	badgeService "gym-badges-api/internal/service/badge"
	"gym-badges-api/models"
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

	forbiddenErrorResponse = models.GenericResponse{
		Code:    fmt.Sprint(http.StatusForbidden),
		Message: http.StatusText(http.StatusForbidden),
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

func (h badgesHandler) GetBadgesByUserID(params op.GetBadgesByUserIDParams) middleware.Responder {

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

func (h badgesHandler) AddBadge(params op.AddBadgeParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("BADGES_HANDLER: Adding badge %d to user %s", params.Input.BadgeID, params.UserID)

	// An user can only edit his own info
	if params.AuthUserID != params.UserID {
		return op.NewAddBadgeUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	err := h.badgeService.AddBadge(params.UserID, int16(params.Input.BadgeID), ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetBadgesByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewAddBadgeNotFound().WithPayload(&notFoundErrorResponse)
		case errors.As(err, &customErrors.Forbidden):
			return op.NewAddBadgeForbidden().WithPayload(&forbiddenErrorResponse)
		default:
			return op.NewAddBadgeInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewAddBadgeOK()
}
