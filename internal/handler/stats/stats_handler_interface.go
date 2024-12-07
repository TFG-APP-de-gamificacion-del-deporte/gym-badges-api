package stats_handler

import (
	"gym-badges-api/restapi/operations/stats"

	"github.com/go-openapi/runtime/middleware"
)

type IStatsHandler interface {
	GetWeightHistory(params stats.GetWeightHistoryByUserIDParams) middleware.Responder
	AddWeight(params stats.AddWeightParams) middleware.Responder

	GetFatHistory(params stats.GetFatHistoryByUserIDParams) middleware.Responder
	AddBodyFat(params stats.AddBodyFatParams) middleware.Responder

	GetStreakCalendar(params stats.GetStreakCalendarByUserIDParams) middleware.Responder
	AddGymAttendance(params stats.AddGymAttendanceParams) middleware.Responder
}
