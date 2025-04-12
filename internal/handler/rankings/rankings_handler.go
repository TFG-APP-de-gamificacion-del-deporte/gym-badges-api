package rankings_handler

import (
	"errors"
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	rankingsService "gym-badges-api/internal/service/rankings"
	"gym-badges-api/models"

	op "gym-badges-api/restapi/operations/rankings"
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

func NewRankingsHandler(rankingsService rankingsService.IRankingsService) IRankingsHandler {
	return &rankingsHandler{
		rankingsService: rankingsService,
	}
}

type rankingsHandler struct {
	rankingsService rankingsService.IRankingsService
}

func (h *rankingsHandler) GetGlobalRanking(params op.GetGlobalRankingParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("RANKINGS_HANDLER: Getting global ranking with user: %s", params.UserID)

	response, err := h.rankingsService.GetGlobalRanking(params.UserID, params.Page, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetGlobalRankingUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetGlobalRankingNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetGlobalRankingInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetGlobalRankingOK().WithPayload(response)
}

func (h *rankingsHandler) GetFriendsRanking(params op.GetFriendsRankingParams) middleware.Responder {

	ctxLog := toolsLogging.BuildLogger(params.HTTPRequest.Context())

	ctxLog.Infof("RANKINGS_HANDLER: Getting friends ranking with user: %s", params.UserID)

	response, err := h.rankingsService.GetFriendsRanking(params.UserID, params.Page, ctxLog)
	if err != nil {
		switch {
		case errors.As(err, &customErrors.Unauthorized):
			return op.NewGetFriendsRankingUnauthorized().WithPayload(&unauthorizedErrorResponse)
		case errors.As(err, &customErrors.NotFound):
			return op.NewGetFriendsRankingNotFound().WithPayload(&notFoundErrorResponse)
		default:
			return op.NewGetFriendsRankingInternalServerError().WithPayload(&internalServerErrorResponse)
		}
	}

	return op.NewGetFriendsRankingOK().WithPayload(response)
}
