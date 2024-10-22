package user_handler

import (
	"gym-badges-api/restapi/operations/user"

	"github.com/go-openapi/runtime/middleware"
)

type IUserHandler interface {
	GetUser(params user.GetUserInfoParams) middleware.Responder
}
