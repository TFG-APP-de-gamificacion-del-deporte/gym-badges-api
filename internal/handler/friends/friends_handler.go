package friends_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	friendsService "gym-badges-api/internal/service/friends"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/friends"
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

func NewFriendsHandler(friendsService friendsService.IFriendsService) IFriendsHandler {
	return &friendsHandler{
		friendsService: friendsService,
	}
}

type friendsHandler struct {
	friendsService friendsService.IFriendsService
}

func (h friendsHandler) GetFriendsByUserID(params op.GetFriendsByUserIDParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("FRIENDS_HANDLER: Getting friends for user: %s", params.UserID)

	response, err := h.friendsService.GetFriendsByUserID(params.UserID, params.Page, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetFriendsByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetFriendsByUserIDNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetFriendsByUserIDInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetFriendsByUserIDOK().WithPayload(response)
}

func (h friendsHandler) AddFriend(params op.AddFriendParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("FRIENDS_HANDLER: Making %s (user) and %s (friend) friends.", params.UserID, params.Input.FriendID)

	// An user can only add friends to himself
	if params.AuthUserID != params.UserID {
		return op.NewAddFriendUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	response, err := h.friendsService.AddFriend(params.UserID, params.Input.FriendID, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewAddFriendUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewAddFriendNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewAddFriendInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewAddFriendOK().WithPayload(response)
}

func (h friendsHandler) DeleteFriend(params op.DeleteFriendParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("FRIENDS_HANDLER: Making %s (user) and %s (friend) no longer friends.", params.UserID, params.Input.FriendID)

	// An user can only delete his own friends
	if params.AuthUserID != params.UserID {
		return op.NewAddFriendUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	err := h.friendsService.DeleteFriend(params.UserID, params.Input.FriendID, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewAddFriendUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewAddFriendNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewAddFriendInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewDeleteFriendOK()
}

func (h friendsHandler) GetFriendRequestsByUserID(params op.GetFriendRequestsByUserIDParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("FRIENDS_HANDLER: Getting friend requests for user: %s", params.UserID)

	// An user can only see his own friend requests
	if params.AuthUserID != params.UserID {
		return op.NewGetFriendRequestsByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	response, err := h.friendsService.GetFriendRequestsByUserID(params.UserID, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetFriendRequestsByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetFriendRequestsByUserIDNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetFriendRequestsByUserIDInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetFriendRequestsByUserIDOK().WithPayload(response)
}
