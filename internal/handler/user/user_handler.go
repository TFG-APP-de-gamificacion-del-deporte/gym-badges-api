package user_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	userService "gym-badges-api/internal/service/user"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/user"
	toolsLogging "gym-badges-api/tools/logging"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

var (
	unauthorizedError customErrors.UnauthorizedError
	conflictError     customErrors.ConflictError
	NotFoundError     customErrors.NotFoundError

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

func NewUserHandler(userService userService.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

type userHandler struct {
	userService userService.IUserService
}

func (h userHandler) GetUser(params op.GetUserInfoParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("USER_HANDLER: Getting info for user: %s", params.UserID)

	response, err := h.userService.GetUser(params.UserID, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &unauthorizedError):
			return op.NewGetUserInfoUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &NotFoundError):
			return op.NewGetUserInfoNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetUserInfoInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetUserInfoOK().WithPayload(response)
}

func (h userHandler) CreateUser(params op.CreateUserParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("USER_HANDLER: Creating user: %s", params.Input.UserID)

	response, err := h.userService.CreateUser(params.Input, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &conflictError):
			return op.NewCreateUserConflict().WithPayload(&models.GenericResponse{
				Code:    fmt.Sprint(http.StatusConflict),
				Message: err.Error(),
			})
		default:
			return op.NewGetUserInfoInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewCreateUserCreated().WithPayload(response)
}

func (h userHandler) EditUserInfo(params op.EditUserInfoParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("USER_HANDLER: Editing user: %s", params.UserID)

	// An user can only edit his own info
	if params.AuthUserID != params.UserID {
		return op.NewEditUserInfoUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	response, err := h.userService.EditUserInfo(params.UserID, params.Input, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &unauthorizedError):
			return op.NewEditUserInfoUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &NotFoundError):
			return op.NewEditUserInfoNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewEditUserInfoInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewEditUserInfoOK().WithPayload(response)
}
