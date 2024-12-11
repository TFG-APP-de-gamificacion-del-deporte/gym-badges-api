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
	"time"

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

	// An user can only add new weights to himself
	if params.AuthUserID != params.UserID {
		return op.NewAddWeightUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

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

	return op.NewAddWeightOK()
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

func (h statsHandler) AddBodyFat(params op.AddBodyFatParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Adding new body fat to user: %s", params.UserID)

	// An user can only add new body fats to himself
	if params.AuthUserID != params.UserID {
		return op.NewAddBodyFatUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	err := h.statsService.AddBodyFat(params.UserID, params.Input.BodyFat, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewAddBodyFatUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewAddBodyFatNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewAddBodyFatInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewAddBodyFatOK()
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

func (h statsHandler) AddGymAttendance(params op.AddGymAttendanceParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Adding a gym attendance to user: %s", params.UserID)

	// An user can only add new gym attendances to himself
	if params.AuthUserID != params.UserID {
		return op.NewAddGymAttendanceUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	// Check date is not in the future
	if time.Time(params.Input.Date).After(time.Now()) {
		return op.NewAddGymAttendanceBadRequest().WithPayload(&models.GenericResponse{
			Code:    fmt.Sprint(http.StatusBadRequest),
			Message: "Date cannot be in the future.",
		})
	}

	err := h.statsService.AddGymAttendance(params.UserID, time.Time(params.Input.Date), ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewAddGymAttendanceUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewAddGymAttendanceNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewAddGymAttendanceInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewAddGymAttendanceOK()
}

func (h statsHandler) DeleteGymAttendance(params op.DeleteGymAttendanceParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("STATS_HANDLER: Deleting a gym attendance to user: %s", params.UserID)

	// An user can only delete gym attendances to himself
	if params.AuthUserID != params.UserID {
		return op.NewDeleteGymAttendanceUnauthorized().WithPayload(&unauthorizedErrorResponse)
	}

	err := h.statsService.DeleteGymAttendance(params.UserID, time.Time(params.Input.Date), ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewDeleteGymAttendanceUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewDeleteGymAttendanceNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewDeleteGymAttendanceInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewDeleteGymAttendanceOK()
}
