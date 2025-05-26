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
	"time"

	"github.com/go-openapi/strfmt"
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

	Context("POST /stats/weight/{user_id}", func() {
		var (
			params op.AddWeightParams
		)

		BeforeEach(func() {
			params = op.NewAddWeightParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
			params.Input = &models.AddWeightRequest{
				Weight: 75.5,
			}
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.AddWeight(params)
			Expect(response).To(BeEquivalentTo(op.NewAddWeightUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking add weight handler cases", func(input Params) {
			mockStatsService.EXPECT().AddWeight(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.AddWeight(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200) - Add weight successfully", Params{
				ExpectedResponse: op.NewAddWeightOK(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to add weight", Params{
				ExpectedResponse: op.NewAddWeightUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewAddWeightNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewAddWeightInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)
	})

	Context("POST /stats/fat/{user_id}", func() {
		var (
			params op.AddBodyFatParams
		)

		BeforeEach(func() {
			params = op.NewAddBodyFatParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
			params.Input = &models.AddBodyFatRequest{
				BodyFat: 15.5,
			}
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.AddBodyFat(params)
			Expect(response).To(BeEquivalentTo(op.NewAddBodyFatUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking add body fat handler cases", func(input Params) {
			mockStatsService.EXPECT().AddBodyFat(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.AddBodyFat(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200) - Add body fat successfully", Params{
				ExpectedResponse: op.NewAddBodyFatOK(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to add body fat", Params{
				ExpectedResponse: op.NewAddBodyFatUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewAddBodyFatNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewAddBodyFatInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)
	})

	Context("POST /stats/attendance/{user_id}", func() {
		var (
			params op.AddGymAttendanceParams
		)

		BeforeEach(func() {
			params = op.NewAddGymAttendanceParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
			params.Input = &models.AddDeleteDayToStreakRequest{
				Date: strfmt.Date(time.Now()),
			}
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.AddGymAttendance(params)
			Expect(response).To(BeEquivalentTo(op.NewAddGymAttendanceUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		It("should return bad request when date is in the future", func() {
			params.Input.Date = strfmt.Date(time.Now().AddDate(0, 0, 1))

			response := handler.AddGymAttendance(params)
			Expect(response).To(BeEquivalentTo(op.NewAddGymAttendanceBadRequest().WithPayload(&models.GenericResponse{
				Code:    "400",
				Message: "Date cannot be in the future.",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking add gym attendance handler cases", func(input Params) {
			mockStatsService.EXPECT().AddGymAttendance(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.AddGymAttendance(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200) - Add gym attendance successfully", Params{
				ExpectedResponse: op.NewAddGymAttendanceOK(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to add gym attendance", Params{
				ExpectedResponse: op.NewAddGymAttendanceUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewAddGymAttendanceNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewAddGymAttendanceInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)
	})

	Context("DELETE /stats/attendance/{user_id}", func() {
		var (
			params op.DeleteGymAttendanceParams
		)

		BeforeEach(func() {
			params = op.NewDeleteGymAttendanceParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
			params.Input = &models.AddDeleteDayToStreakRequest{
				Date: strfmt.Date(time.Now()),
			}
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.DeleteGymAttendance(params)
			Expect(response).To(BeEquivalentTo(op.NewDeleteGymAttendanceUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking delete gym attendance handler cases", func(input Params) {
			mockStatsService.EXPECT().DeleteGymAttendance(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.DeleteGymAttendance(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200) - Delete gym attendance successfully", Params{
				ExpectedResponse: op.NewDeleteGymAttendanceOK(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to delete gym attendance", Params{
				ExpectedResponse: op.NewDeleteGymAttendanceUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewDeleteGymAttendanceNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewDeleteGymAttendanceInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)
	})

})
