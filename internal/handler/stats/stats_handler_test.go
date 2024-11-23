package stats_handler

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/mocks/service"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/stats"
	toolsTesting "gym-badges-api/tools/testing"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHandlerStatsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "HANDLER: Stats Test Suite")
}

var _ = Describe("HANDLER: Stats Test Suite", func() {

	var (
		mockCtrl         *gomock.Controller
		mockStatsService *service.MockIStatsService
		handler          IStatsHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockStatsService = service.NewMockIStatsService(mockCtrl)

		handler = NewStatsHandler(mockStatsService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()

	})

	Context("GET /stats/weight/{user_id}", func() {

		var (
			params op.GetWeightHistoryByUserIDParams
		)

		BeforeEach(func() {
			params = op.NewGetWeightHistoryByUserIDParams()
			params.HTTPRequest = new(http.Request)
			params.Months = 3
			params.UserID = "admin"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.MeasurementHistoryResponse
			ServiceError     error
		}

		DescribeTable("Checking get weight history handler cases", func(input Params) {

			mockStatsService.EXPECT().GetWeightHistory(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetWeightHistory(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetWeightHistoryByUserIDOK().WithPayload(&models.MeasurementHistoryResponse{
					Days: []*models.MeasurementPerDay{
						{
							Date:  "2024-11-01",
							Value: 78.5,
						},
						{
							Date:  "2024-11-07",
							Value: 79,
						},
						{
							Date:  "2024-11-14",
							Value: 80,
						},
					},
				}),
				ServiceResponse: &models.MeasurementHistoryResponse{
					Days: []*models.MeasurementPerDay{
						{
							Date:  "2024-11-01",
							Value: 78.5,
						},
						{
							Date:  "2024-11-07",
							Value: 79,
						},
						{
							Date:  "2024-11-14",
							Value: 80,
						},
					},
				},
				ServiceError: nil,
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewGetWeightHistoryByUserIDNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetWeightHistoryByUserIDUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetWeightHistoryByUserIDInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

	Context("GET /stats/fat/{user_id}", func() {

		var (
			params op.GetFatHistoryByUserIDParams
		)

		BeforeEach(func() {
			params = op.NewGetFatHistoryByUserIDParams()
			params.HTTPRequest = new(http.Request)
			params.Months = 3
			params.UserID = "admin"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.MeasurementHistoryResponse
			ServiceError     error
		}

		DescribeTable("Checking get fat history handler cases", func(input Params) {

			mockStatsService.EXPECT().GetFatHistory(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetFatHistory(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetFatHistoryByUserIDOK().WithPayload(&models.MeasurementHistoryResponse{
					Days: []*models.MeasurementPerDay{
						{
							Date:  "2024-11-01",
							Value: 78.5,
						},
						{
							Date:  "2024-11-07",
							Value: 79,
						},
						{
							Date:  "2024-11-14",
							Value: 80,
						},
					},
				}),
				ServiceResponse: &models.MeasurementHistoryResponse{
					Days: []*models.MeasurementPerDay{
						{
							Date:  "2024-11-01",
							Value: 78.5,
						},
						{
							Date:  "2024-11-07",
							Value: 79,
						},
						{
							Date:  "2024-11-14",
							Value: 80,
						},
					},
				},
				ServiceError: nil,
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewGetFatHistoryByUserIDNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetFatHistoryByUserIDUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetFatHistoryByUserIDInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

	Context("GET /stats/streak/{user_id}", func() {

		var (
			params op.GetStreakCalendarByUserIDParams
		)

		BeforeEach(func() {
			params = op.NewGetStreakCalendarByUserIDParams()
			params.HTTPRequest = new(http.Request)
			params.Month = 11
			params.Year = 2024
			params.UserID = "admin"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.StreakCalendarResponse
			ServiceError     error
		}

		DescribeTable("Checking get streak calendar handler cases", func(input Params) {

			mockStatsService.EXPECT().GetStreakCalendarByYearAndMonth(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetStreakCalendar(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetStreakCalendarByUserIDOK().WithPayload(&models.StreakCalendarResponse{
					Days: []string{
						"2024-11-01",
						"2024-11-02",
						"2024-11-05",
						"2024-11-07",
					},
					Streak:     77,
					WeeklyGoal: 3,
				}),
				ServiceResponse: &models.StreakCalendarResponse{
					Days: []string{
						"2024-11-01",
						"2024-11-02",
						"2024-11-05",
						"2024-11-07",
					},
					Streak:     77,
					WeeklyGoal: 3,
				},
				ServiceError: nil,
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewGetStreakCalendarByUserIDNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetStreakCalendarByUserIDUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetStreakCalendarByUserIDInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

})
