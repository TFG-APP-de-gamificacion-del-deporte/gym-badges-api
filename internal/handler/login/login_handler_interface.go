package login_handler

import (
	"gym-badges-api/restapi/operations/login"
	"gym-badges-api/restapi/operations/login_with_token"

	"github.com/go-openapi/runtime/middleware"
)

type ILoginHandler interface {
	Login(params login.LoginParams) middleware.Responder

	LoginWithToken(params login_with_token.LoginWithTokenParams) middleware.Responder
}
