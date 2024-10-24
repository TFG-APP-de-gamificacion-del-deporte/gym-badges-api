package login_handler

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	"gym-badges-api/mocks/service"
	"gym-badges-api/models"
	op "gym-badges-api/restapi/operations/login"
	toolsTesting "gym-badges-api/tools/testing"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func TestHandlerLoginSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "HANDLER: Login Test Suite")
}

var _ = Describe("HANDLER: Login Test Suite", func() {

	var (
		mockCtrl         *gomock.Controller
		mockLoginService *service.MockILoginService
		handler          ILoginHandler
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockLoginService = service.NewMockILoginService(mockCtrl)

		handler = NewLoginHandler(mockLoginService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()

	})

	Context("POST /Login", func() {

		var (
			params op.LoginParams
		)

		BeforeEach(func() {
			params = op.NewLoginParams()
			params.HTTPRequest = new(http.Request)
			params.Input = new(models.LoginRequestBody)
		})

		type Params struct {
			ExpectedResponse any
			ServiceResponse  *models.LoginResponse
			ServiceError     error
		}

		DescribeTable("Checking login handler cases", func(input Params) {

			mockLoginService.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(input.ServiceResponse, input.ServiceError)

			response := handler.Login(params)
			Expect(response).To(BeEquivalentTo(input.ExpectedResponse))
		},
			Entry("CASE: Success Response (200)", Params{
				ExpectedResponse: op.NewLoginOK().WithPayload(&models.LoginResponse{
					Token: "<PASSWORD>",
				}),
				ServiceResponse: &models.LoginResponse{
					Token: "<PASSWORD>",
				},
				ServiceError: nil,
			}),
			Entry("CASE: Unauthorized Error Response (401)", Params{
				ExpectedResponse: op.NewLoginUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "401",
					Message: "Unauthorized",
				}),
				ServiceResponse: nil,
				ServiceError:    customErrors.BuildUnauthorizedError("unauthorized"),
			}),
			Entry("CASE: Internal Server Error Response (500)", Params{
				ExpectedResponse: op.NewLoginUnauthorized().WithPayload(&models.GenericResponse{
					Code:    "500",
					Message: "Internal Server Error",
				}),
				ServiceResponse: nil,
				ServiceError:    errors.New("panic"),
			}),
		)

	})

})
