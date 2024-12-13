package friends_handler

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/mocks/service"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/friends"
	toolsTesting "gym-badges-api/tools/testing"
	"gym-badges-api/tools/utils"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHandlerFriendsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "HANDLER: Friends Test Suite")
}

var _ = Describe("HANDLER: Friends Test Suite", func() {

	var (
		mockCtrl           *gomock.Controller
		mockFriendsService *service.MockIFriendsService
		handler            IFriendsHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockFriendsService = service.NewMockIFriendsService(mockCtrl)

		handler = NewFriendsHandler(mockFriendsService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()

	})

	Context("GET /friends/{user_id}", func() {

		var (
			params op.GetFriendsByUserIDParams
		)

		BeforeEach(func() {
			params = op.NewGetFriendsByUserIDParams()
			params.HTTPRequest = new(http.Request)
			params.Page = 1
			params.UserID = "admin"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.FriendsResponse
			ServiceError     error
		}

		DescribeTable("Checking get friends by user_id handler cases", func(input Params) {

			mockFriendsService.EXPECT().GetFriendsByUserID(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetFriendsByUserID(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetFriendsByUserIDOK().WithPayload(&models.FriendsResponse{
					Friends: []*models.FriendInfo{
						{
							Fat:    utils.NewFloat32(5.5),
							Image:  "/friend1.jpg",
							Level:  10,
							Name:   "Friend 1",
							Streak: 10,
							TopFeats: []*models.Feat{
								{
									Description: "Description 1",
									Image:       "/feat_1.jpg",
									Name:        "Feat 1",
								},
							},
							User:   "friend1",
							Weight: utils.NewFloat32(80.5),
						},
						{
							Fat:    utils.NewFloat32(0.5),
							Image:  "/friend2.jpg",
							Level:  20,
							Name:   "Friend 2",
							Streak: 20,
							TopFeats: []*models.Feat{
								{
									Description: "Description 2",
									Image:       "/feat_2.jpg",
									Name:        "Feat 2",
								},
							},
							User:   "friend2",
							Weight: utils.NewFloat32(70.5),
						},
					},
				}),
				ServiceResponse: &models.FriendsResponse{
					Friends: []*models.FriendInfo{
						{
							Fat:    utils.NewFloat32(5.5),
							Image:  "/friend1.jpg",
							Level:  10,
							Name:   "Friend 1",
							Streak: 10,
							TopFeats: []*models.Feat{
								{
									Description: "Description 1",
									Image:       "/feat_1.jpg",
									Name:        "Feat 1",
								},
							},
							User:   "friend1",
							Weight: utils.NewFloat32(80.5),
						},
						{
							Fat:    utils.NewFloat32(0.5),
							Image:  "/friend2.jpg",
							Level:  20,
							Name:   "Friend 2",
							Streak: 20,
							TopFeats: []*models.Feat{
								{
									Description: "Description 2",
									Image:       "/feat_2.jpg",
									Name:        "Feat 2",
								},
							},
							User:   "friend2",
							Weight: utils.NewFloat32(70.5),
						},
					},
				},
				ServiceError: nil,
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewGetFriendsByUserIDNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetFriendsByUserIDUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetFriendsByUserIDInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

})
