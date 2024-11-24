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
			ServiceResponse  *models.BadgesByUserResponse
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
				ExpectedResponse: op.NewGetBadgesByUserIDOK().WithPayload(&models.BadgesByUserResponse{
					Badges: map[string]models.Badge{
						"1": {
							Achieved:    false,
							Description: "Badge 1",
							Image:       "/badge_1.jpg",
							Name:        "Badge 1",
							Parent:      "",
						},
						"2": {
							Achieved:    false,
							Description: "Badge 2",
							Image:       "/badge_2.jpg",
							Name:        "Badge 2",
							Parent:      "1",
						},
						"3": {
							Achieved:    false,
							Description: "Badge 3",
							Image:       "/badge_3.jpg",
							Name:        "Badge 3",
							Parent:      "1",
						},
					},
				}),
				ServiceResponse: &models.BadgesByUserResponse{
					Badges: map[string]models.Badge{
						"1": {
							Achieved:    false,
							Description: "Badge 1",
							Image:       "/badge_1.jpg",
							Name:        "Badge 1",
							Parent:      "",
						},
						"2": {
							Achieved:    false,
							Description: "Badge 2",
							Image:       "/badge_2.jpg",
							Name:        "Badge 2",
							Parent:      "1",
						},
						"3": {
							Achieved:    false,
							Description: "Badge 3",
							Image:       "/badge_3.jpg",
							Name:        "Badge 3",
							Parent:      "1",
						},
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

})
