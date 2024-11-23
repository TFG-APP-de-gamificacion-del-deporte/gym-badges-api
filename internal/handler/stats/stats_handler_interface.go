package stats_handler

import (
	"gym-badges-api/restapi/operations/stats"

	"github.com/go-openapi/runtime/middleware"
)

type IStatsHandler interface {
	GetWeightHistory(params stats.GetWeightHistoryByUserIDParams) middleware.Responder
}
