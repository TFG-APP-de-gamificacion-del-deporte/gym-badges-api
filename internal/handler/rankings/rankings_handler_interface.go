package rankings_handler

import (
	"gym-badges-api/restapi/operations/rankings"

	"github.com/go-openapi/runtime/middleware"
)

type IRankingsHandler interface {
	GetGlobalRanking(params rankings.GetGlobalRankingParams) middleware.Responder
}
