package user_handler

import (
	userService "gym-badges-api/internal/service/user"
	"gym-badges-api/restapi/operations/user"

	"github.com/go-openapi/runtime/middleware"
)

func NewUserHandler(userService userService.IUserService) IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

type userHandler struct {
	userService userService.IUserService
}

func (h userHandler) GetUser(params user.GetUserInfoParams) middleware.Responder {
	return nil
}
