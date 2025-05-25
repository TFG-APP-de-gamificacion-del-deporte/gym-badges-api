package rankings_handler

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/mocks/service"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/rankings"
	toolsTesting "gym-badges-api/tools/testing"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHandlerRankingsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "HANDLER: Rankings Test Suite")
}

var _ = Describe("HANDLER: Rankings Test Suite", func() {

	var (
		mockCtrl            *gomock.Controller
		mockRankingsService *service.MockIRankingsService
		handler             IRankingsHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRankingsService = service.NewMockIRankingsService(mockCtrl)

		handler = NewRankingsHandler(mockRankingsService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("GET /rankings/global/{user_id}", func() {
		var (
			params op.GetGlobalRankingParams
		)

		BeforeEach(func() {
			params = op.NewGetGlobalRankingParams()
			params.HTTPRequest = new(http.Request)
			params.Page = 1
			params.UserID = "test_user"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.GetRankingResponse
			ServiceError     error
		}

		DescribeTable("Checking get global ranking handler cases", func(input Params) {
			mockRankingsService.EXPECT().GetGlobalRanking(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetGlobalRanking(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetGlobalRankingOK().WithPayload(&models.GetRankingResponse{
					Ranking: []*models.RakingUser{
						{
							UserID: "user1",
							Name:   "John Doe",
							Image:  []byte("image1"),
							Level:  10,
							Rank:   1,
							Streak: 5,
						},
						{
							UserID: "user2",
							Name:   "Jane Doe",
							Image:  []byte("image2"),
							Level:  8,
							Rank:   2,
							Streak: 3,
						},
					},
					Yourself: &models.RakingUser{
						UserID: "test_user",
						Name:   "Test User",
						Image:  []byte("image3"),
						Level:  5,
						Rank:   15,
						Streak: 2,
					},
				}),
				ServiceResponse: &models.GetRankingResponse{
					Ranking: []*models.RakingUser{
						{
							UserID: "user1",
							Name:   "John Doe",
							Image:  []byte("image1"),
							Level:  10,
							Rank:   1,
							Streak: 5,
						},
						{
							UserID: "user2",
							Name:   "Jane Doe",
							Image:  []byte("image2"),
							Level:  8,
							Rank:   2,
							Streak: 3,
						},
					},
					Yourself: &models.RakingUser{
						UserID: "test_user",
						Name:   "Test User",
						Image:  []byte("image3"),
						Level:  5,
						Rank:   15,
						Streak: 2,
					},
				},
				ServiceError: nil,
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewGetGlobalRankingNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetGlobalRankingUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetGlobalRankingInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)
	})

	Context("GET /rankings/friends/{user_id}", func() {
		var (
			params op.GetFriendsRankingParams
		)

		BeforeEach(func() {
			params = op.NewGetFriendsRankingParams()
			params.HTTPRequest = new(http.Request)
			params.Page = 1
			params.UserID = "test_user"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.GetRankingResponse
			ServiceError     error
		}

		DescribeTable("Checking get friends ranking handler cases", func(input Params) {
			mockRankingsService.EXPECT().GetFriendsRanking(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetFriendsRanking(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewGetFriendsRankingOK().WithPayload(&models.GetRankingResponse{
					Ranking: []*models.RakingUser{
						{
							UserID: "friend1",
							Name:   "Friend One",
							Image:  []byte("image1"),
							Level:  10,
							Rank:   1,
							Streak: 5,
						},
						{
							UserID: "friend2",
							Name:   "Friend Two",
							Image:  []byte("image2"),
							Level:  8,
							Rank:   2,
							Streak: 3,
						},
					},
					Yourself: &models.RakingUser{
						UserID: "test_user",
						Name:   "Test User",
						Image:  []byte("image3"),
						Level:  5,
						Rank:   15,
						Streak: 2,
					},
				}),
				ServiceResponse: &models.GetRankingResponse{
					Ranking: []*models.RakingUser{
						{
							UserID: "friend1",
							Name:   "Friend One",
							Image:  []byte("image1"),
							Level:  10,
							Rank:   1,
							Streak: 5,
						},
						{
							UserID: "friend2",
							Name:   "Friend Two",
							Image:  []byte("image2"),
							Level:  8,
							Rank:   2,
							Streak: 3,
						},
					},
					Yourself: &models.RakingUser{
						UserID: "test_user",
						Name:   "Test User",
						Image:  []byte("image3"),
						Level:  5,
						Rank:   15,
						Streak: 2,
					},
				},
				ServiceError: nil,
			}),
			Entry("CASE: Not Found Error Response (404)", Params{
				ExpectedResponse: op.NewGetFriendsRankingNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewGetFriendsRankingUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewGetFriendsRankingInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)
	})
})
