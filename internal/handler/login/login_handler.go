package login_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	loginService "gym-badges-api/internal/service/login"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/login"
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

func NewLoginHandler(loginService loginService.ILoginService) ILoginHandler {
	return &loginHandler{
		loginService: loginService,
	}
}

type loginHandler struct {
	loginService loginService.ILoginService
}

func (receiver loginHandler) Login(params op.LoginParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("LOGIN_HANDLER: Login for user: %s", params.Input.User)

	response, err := receiver.loginService.Login(params.Input.User, params.Input.Password, ctxLog)
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
