package login_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	loginService "gym-badges-api/internal/service/login"
	sessionService "gym-badges-api/internal/service/session"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/login"
	"gym-badges-api/restapi/operations/login_with_token"
	toolsLogging "gym-badges-api/tools/logging"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

var (
	unauthorizedError customErrors.UnauthorizedError

	unauthorizedErrorResponse = models.GenericResponse{
		Code:    fmt.Sprint(http.StatusUnauthorized),
		Message: http.StatusText(http.StatusUnauthorized),
	}

	internalServerErrorResponse = models.GenericResponse{
		Code:    fmt.Sprint(http.StatusInternalServerError),
		Message: http.StatusText(http.StatusInternalServerError),
	}
)

func NewLoginHandler(loginService loginService.ILoginService, sessionService sessionService.ISessionService) ILoginHandler {
	return &loginHandler{
		loginService:   loginService,
		sessionService: sessionService,
	}
}

type loginHandler struct {
	loginService   loginService.ILoginService
	sessionService sessionService.ISessionService
}

func (h loginHandler) Login(params op.LoginParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("LOGIN_HANDLER: Login for user: %s", params.Input.User)

	response, err := h.loginService.Login(params.Input.User, params.Input.Password, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &unauthorizedError):
			return op.NewLoginUnauthorized().WithPayload(&unauthorizedErrorResponse)
		default:
			return op.NewLoginInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewLoginOK().WithPayload(response)

}

func (h loginHandler) LoginWithToken(params login_with_token.LoginWithTokenParams) middleware.Responder {

	username, err := h.sessionService.GetUserFromToken(params.Token)

	if err != nil {
		return login_with_token.NewLoginWithTokenUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	response := models.LoginWithTokenResponse{
		Username: username,
	}

	return login_with_token.NewLoginWithTokenOK().WithPayload(&response)
}
