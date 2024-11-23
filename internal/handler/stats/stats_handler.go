package stats_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	statsService "gym-badges-api/internal/service/stats"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/stats"
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

func NewStatsHandler(statsService statsService.IStatsService) IStatsHandler {
	return &statsHandler{
		statsService: statsService,
	}
}

type statsHandler struct {
	statsService statsService.IStatsService
}

func (h statsHandler) GetWeightHistory(params op.GetWeightHistoryByUserIDParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Getting weight history for user: %s", params.UserID)

	response, err := h.statsService.GetWeightHistory(params.UserID, params.Months, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetWeightHistoryByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetWeightHistoryByUserIDNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetWeightHistoryByUserIDInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetWeightHistoryByUserIDOK().WithPayload(response)
}
