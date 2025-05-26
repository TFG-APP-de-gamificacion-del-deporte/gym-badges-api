package badge_handler

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/mocks/service"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/badges"
	toolsTesting "gym-badges-api/tools/testing"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHandlerBadgeSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "HANDLER: Badge Test Suite")
}

var _ = Describe("HANDLER: Badge Test Suite", func() {

	var (
		mockCtrl         *gomock.Controller
		mockBadgeService *service.MockIBadgeService
		handler          IBadgeHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockBadgeService = service.NewMockIBadgeService(mockCtrl)

		handler = NewBadgeHandler(mockBadgeService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()

	})

	Context("GET /badges/{user_id}", func() {

		var (
			params op.GetBadgesByUserIDParams
		)

		BeforeEach(func() {
			params = op.NewGetBadgesByUserIDParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "admin"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  models.BadgesByUserResponse
			ServiceError     error
		}

		DescribeTable("Checking get badges by user_id handler cases", func(input Params) {

			mockBadgeService.EXPECT().GetBadgesByUserID(gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetBadgesByUserID(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetBadgesByUserIDOK().WithPayload(models.BadgesByUserResponse{
					{
						Achieved:    false,
						Description: "Badge 1",
						ID:          1,
						Image:       "/badge_1.jpg",
						Name:        "Badge 1",
					},
					{
						Achieved:    false,
						Description: "Badge 2",
						ID:          2,
						Image:       "/badge_2.jpg",
						Name:        "Badge 2",
					},
					{
						Achieved:    false,
						Description: "Badge 3",
						ID:          3,
						Image:       "/badge_3.jpg",
						Name:        "Badge 3",
					},
				}),
				ServiceResponse: models.BadgesByUserResponse{
					{
						Achieved:    false,
						Description: "Badge 1",
						ID:          1,
						Image:       "/badge_1.jpg",
						Name:        "Badge 1",
					},
					{
						Achieved:    false,
						Description: "Badge 2",
						ID:          2,
						Image:       "/badge_2.jpg",
						Name:        "Badge 2",
					},
					{
						Achieved:    false,
						Description: "Badge 3",
						ID:          3,
						Image:       "/badge_3.jpg",
						Name:        "Badge 3",
					},
				},
				ServiceError: nil,
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetBadgesByUserIDUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetBadgesByUserIDInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

	Context("POST /badges/{user_id}", func() {

		var (
			params op.AddBadgeParams
		)

		BeforeEach(func() {
			params = op.NewAddBadgeParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "admin"
			params.AuthUserID = "admin"
			params.Input = &models.AddDeleteBadgeRequest{
				BadgeID: 1,
			}
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking add badges by user_id handler cases", func(input Params) {

			mockBadgeService.EXPECT().AddBadge(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.AddBadge(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewAddBadgeOK(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewAddBadgeUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Forbidden Error Response (403)", Params{
				ExpectedResponse: op.NewAddBadgeForbidden().WithPayload(&models.GenericResponse{
					Code:    "403",
					Message: "Forbidden",
				}),
				ServiceError: customErrors.BuildForbiddenError("forbidden"),
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewAddBadgeNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("not found"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewAddBadgeInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)

		It("CASE: Unauthorized Error Response (401) - Bad user id", func() {

			params.AuthUserID = "other"

			expectedResponse := op.NewAddBadgeUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})
			response := handler.AddBadge(params)
			Expect(response).To(BeEquivalentTo(expectedResponse))
		})

	})

	Context("DELETE /badges/{user_id}", func() {

		var (
			params op.DeleteBadgeParams
		)

		BeforeEach(func() {
			params = op.NewDeleteBadgeParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "admin"
			params.AuthUserID = "admin"
			params.Input = &models.AddDeleteBadgeRequest{
				BadgeID: 1,
			}
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking add badges by user_id handler cases", func(input Params) {

			mockBadgeService.EXPECT().DeleteBadge(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.DeleteBadge(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewDeleteBadgeOK(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewDeleteBadgeUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Forbidden Error Response (403)", Params{
				ExpectedResponse: op.NewDeleteBadgeForbidden().WithPayload(&models.GenericResponse{
					Code:    "403",
					Message: "Forbidden",
				}),
				ServiceError: customErrors.BuildForbiddenError("forbidden"),
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewDeleteBadgeNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("not found"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewDeleteBadgeInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)

		It("CASE: Unauthorized Error Response (401) - Bad user id", func() {

			params.AuthUserID = "other"

			expectedResponse := op.NewDeleteBadgeUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})
			response := handler.DeleteBadge(params)
			Expect(response).To(BeEquivalentTo(expectedResponse))
		})

	})

})
