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

func (h statsHandler) AddWeight(params op.AddWeightParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Adding new weight to user: %s", params.UserID)

	err := h.statsService.AddWeight(params.UserID, params.Input.Weight, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewAddWeightUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewAddWeightNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewAddWeightInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return nil
}

func (h statsHandler) GetFatHistory(params op.GetFatHistoryByUserIDParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Getting fat history for user: %s", params.UserID)

	response, err := h.statsService.GetFatHistory(params.UserID, params.Months, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetFatHistoryByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetFatHistoryByUserIDNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetFatHistoryByUserIDInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetFatHistoryByUserIDOK().WithPayload(response)
}

func (h statsHandler) GetStreakCalendar(params op.GetStreakCalendarByUserIDParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Getting streak calendar for user: %s in year: %d month: %d", params.UserID,
		params.Year, params.Month)

	response, err := h.statsService.GetStreakCalendarByYearAndMonth(params.UserID, params.Year, params.Month, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetStreakCalendarByUserIDUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetStreakCalendarByUserIDNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetStreakCalendarByUserIDInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetStreakCalendarByUserIDOK().WithPayload(response)
}
