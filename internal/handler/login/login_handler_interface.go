package login_handler

import (
	"gym-badges-api/restapi/operations/login"

	"github.com/go-openapi/runtime/middleware"
)

type ILoginHandler interface {
	Login(params login.LoginParams) middleware.Responder
}
