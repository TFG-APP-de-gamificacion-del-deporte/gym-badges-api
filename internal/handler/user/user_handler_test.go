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

})
