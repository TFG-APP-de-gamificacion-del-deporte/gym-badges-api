package user_handler

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/mocks/service"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/user"
	toolsTesting "gym-badges-api/tools/testing"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func float32Ptr(v float32) *float32 {
	return &v
}

func TestHandlerUserSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "HANDLER: User Test Suite")
}

var _ = Describe("HANDLER: User Test Suite", func() {

	var (
		mockCtrl        *gomock.Controller
		mockUserService *service.MockIUserService
		handler         IUserHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserService = service.NewMockIUserService(mockCtrl)

		handler = NewUserHandler(mockUserService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()

	})

	Context("POST /user", func() {

		var (
			params op.CreateUserParams
		)

		BeforeEach(func() {
			params = op.NewCreateUserParams()
			params.HTTPRequest = new(http.Request)
			params.Input = new(models.CreateUserRequest)
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.LoginResponse
			ServiceError     error
		}

		DescribeTable("Checking user creation handler cases", func(input Params) {

			mockUserService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.CreateUser(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewCreateUserCreated().WithPayload(&models.LoginResponse{
					Token: "<PASSWORD>",
				}),
				ServiceResponse: &models.LoginResponse{
					Token: "<PASSWORD>",
				},
				ServiceError: nil,
			}),
			Entry("CASE: Conflict Error Response (409)", Params{
				ExpectedResponse: op.NewCreateUserConflict().WithPayload(&models.GenericResponse{
					Code:    "409",
					Message: "user already exists",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildConflictError("user already exists"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewCreateUserInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

	Context("POST /user/user_id", func() {

		var (
			params op.EditUserInfoParams
		)

		BeforeEach(func() {
			params = op.NewEditUserInfoParams()
			params.HTTPRequest = new(http.Request)
			params.Input = new(models.EditUserInfoRequest)
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.EditUserInfo(params)
			Expect(response).To(BeEquivalentTo(op.NewEditUserInfoUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.GetUserInfoResponse
			ServiceError     error
		}

		DescribeTable("Checking user edition handler cases", func(input Params) {
			mockUserService.EXPECT().EditUserInfo(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.EditUserInfo(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200) - Update user info successfully", Params{
				ExpectedResponse: op.NewEditUserInfoOK().WithPayload(&models.GetUserInfoResponse{
					UserID:      "test_user",
					BodyFat:     float32Ptr(15.5),
					CurrentWeek: []bool{true, false, true, false, true, false, false},
					Experience:  1000,
					Height:      175,
					Image:       []byte("test_image"),
					IsFriend:    true,
					Name:        "Updated User",
					Preferences: []*models.Preference{
						{PreferenceID: 1, On: true},
						{PreferenceID: 2, On: false},
					},
					Sex:    "masculine",
					Streak: 5,
					TopFeats: []*models.Feat{
						{ID: 1, Name: "Feat 1", Description: "Description 1", Image: "image1"},
						{ID: 2, Name: "Feat 2", Description: "Description 2", Image: "image2"},
					},
					TotalFriends: 10,
					WeeklyGoal:   4,
					Weight:       float32Ptr(75.5),
				}),
				ServiceResponse: &models.GetUserInfoResponse{
					UserID:      "test_user",
					BodyFat:     float32Ptr(15.5),
					CurrentWeek: []bool{true, false, true, false, true, false, false},
					Experience:  1000,
					Height:      175,
					Image:       []byte("test_image"),
					IsFriend:    true,
					Name:        "Updated User",
					Preferences: []*models.Preference{
						{PreferenceID: 1, On: true},
						{PreferenceID: 2, On: false},
					},
					Sex:    "masculine",
					Streak: 5,
					TopFeats: []*models.Feat{
						{ID: 1, Name: "Feat 1", Description: "Description 1", Image: "image1"},
						{ID: 2, Name: "Feat 2", Description: "Description 2", Image: "image2"},
					},
					TotalFriends: 10,
					WeeklyGoal:   4,
					Weight:       float32Ptr(75.5),
				},
				ServiceError: nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to edit info", Params{
				ExpectedResponse: op.NewEditUserInfoUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewEditUserInfoNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewEditUserInfoInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

	Context("GET /user/user_id", func() {
		var (
			params op.GetUserInfoParams
		)

		BeforeEach(func() {
			params = op.NewGetUserInfoParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.GetUserInfoResponse
			ServiceError     error
		}

		DescribeTable("Checking get user handler cases", func(input Params) {
			mockUserService.EXPECT().GetUser(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.GetUser(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200) - Get own user info", Params{
				ExpectedResponse: op.NewGetUserInfoOK().WithPayload(&models.GetUserInfoResponse{
					UserID:      "test_user",
					BodyFat:     float32Ptr(15.5),
					CurrentWeek: []bool{true, false, true, false, true, false, false},
					Experience:  1000,
					Height:      175,
					Image:       []byte("test_image"),
					IsFriend:    true,
					Name:        "Test User",
					Preferences: []*models.Preference{
						{PreferenceID: 1, On: true},
						{PreferenceID: 2, On: false},
					},
					Sex:    "masculine",
					Streak: 5,
					TopFeats: []*models.Feat{
						{ID: 1, Name: "Feat 1", Description: "Description 1", Image: "image1"},
						{ID: 2, Name: "Feat 2", Description: "Description 2", Image: "image2"},
					},
					TotalFriends: 10,
					WeeklyGoal:   3,
					Weight:       float32Ptr(75.5),
				}),
				ServiceResponse: &models.GetUserInfoResponse{
					UserID:      "test_user",
					BodyFat:     float32Ptr(15.5),
					CurrentWeek: []bool{true, false, true, false, true, false, false},
					Experience:  1000,
					Height:      175,
					Image:       []byte("test_image"),
					IsFriend:    true,
					Name:        "Test User",
					Preferences: []*models.Preference{
						{PreferenceID: 1, On: true},
						{PreferenceID: 2, On: false},
					},
					Sex:    "masculine",
					Streak: 5,
					TopFeats: []*models.Feat{
						{ID: 1, Name: "Feat 1", Description: "Description 1", Image: "image1"},
						{ID: 2, Name: "Feat 2", Description: "Description 2", Image: "image2"},
					},
					TotalFriends: 10,
					WeeklyGoal:   3,
					Weight:       float32Ptr(75.5),
				},
				ServiceError: nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to view info", Params{
				ExpectedResponse: op.NewGetUserInfoUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewGetUserInfoNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewGetUserInfoInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)
	})

	Context("PUT /user/user_id/preferences", func() {
		var (
			params op.EditUserPreferencesParams
		)

		BeforeEach(func() {
			params = op.NewEditUserPreferencesParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
			params.Input = &models.EditUserPreferenceRequest{
				PrivateAccount:   true,
				HideWeightAndFat: false,
			}
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.EditUserPreferences(params)
			Expect(response).To(BeEquivalentTo(op.NewEditUserPreferencesUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking edit user preferences handler cases", func(input Params) {
			mockUserService.EXPECT().EditUserPreferences(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.EditUserPreferences(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (204) - Update preferences successfully", Params{
				ExpectedResponse: op.NewEditUserPreferencesNoContent(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to edit preferences", Params{
				ExpectedResponse: op.NewEditUserPreferencesUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewEditUserPreferencesNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewEditUserPreferencesInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)
	})

	Context("PUT /user/user_id/top-feats", func() {
		var (
			params op.EditUserTopFeatsParams
		)

		BeforeEach(func() {
			params = op.NewEditUserTopFeatsParams()
			params.HTTPRequest = new(http.Request)
			params.UserID = "test_user"
			params.AuthUserID = "test_user"
			params.Input = []int16{1, 2, 3}
		})

		It("should return unauthorized when auth user id does not match user id", func() {
			params.UserID = "test_user"
			params.AuthUserID = "different_user"

			response := handler.EditUserTopFeats(params)
			Expect(response).To(BeEquivalentTo(op.NewEditUserTopFeatsUnauthorized().WithPayload(&models.GenericResponse{
				Code:    "401",
				Message: "Unauthorized",
			})))
		})

		type Params struct {
			ExpectedResponse any
			ServiceError     error
		}

		DescribeTable("Checking edit user top feats handler cases", func(input Params) {
			mockUserService.EXPECT().EditUserTopFeats(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceError)

			response := handler.EditUserTopFeats(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (204) - Update top feats successfully", Params{
				ExpectedResponse: op.NewEditUserTopFeatsNoContent(),
				ServiceError:     nil,
			}),
			Entry("CASE: Unauthorized Error Response (401) - User not authorized to edit top feats", Params{
				ExpectedResponse: op.NewEditUserTopFeatsUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceError: customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Not Found Error Response (404) - User does not exist", Params{
				ExpectedResponse: op.NewEditUserTopFeatsNotFound().WithPayload(&models.GenericResponse{
					Code:    "404",
					Message: "Not Found",
				}),
				ServiceError: customErrors.BuildNotFoundError("user not found"),
			}),
			Entry("CASE: Internal Server Error Response (500) - Unexpected error", Params{
				ExpectedResponse: op.NewEditUserTopFeatsInternalServerError().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceError: errors.New("panic"),
			}),
		)
	})

})
