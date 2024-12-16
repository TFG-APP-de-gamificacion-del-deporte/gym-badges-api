package rankings_handler

import (
	// "errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	userService "gym-badges-api/internal/service/user"
	"gym-badges-api/models"
	"gym-badges-api/restapi/operations/rankings"

	// op "gym-badges-api/restapi/operations/user"
	// toolsLogging "gym-badges-api/tools/logging"
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

func NewRankingsHandler(userService userService.IUserService) IRankingsHandler {
	return &rankingsHandler{
		// userService: userService,
	}
}

type rankingsHandler struct {
	// userService userService.IUserService
}

func (r *rankingsHandler) GetGlobalRanking(params rankings.GetGlobalRankingParams) middleware.Responder {
	panic("unimplemented")
}
