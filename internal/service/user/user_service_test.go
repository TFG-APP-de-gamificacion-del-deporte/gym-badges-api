package user_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	mockService "gym-badges-api/mocks/service"
	"gym-badges-api/models"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestServiceUserSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: User Test Suite")
}

var _ = Describe("SERVICE: User Test Suite", func() {

	var (
		mockCtrl           *gomock.Controller
		mockUserDAO        *mockDAO.MockIUserDAO
		mockSessionService *mockService.MockISessionService
		service            IUserService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		mockSessionService = mockService.NewMockISessionService(mockCtrl)
		service = NewUserService(mockUserDAO, mockSessionService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("Create User", func() {

		var (
			ctxLogger *log.Entry

			request models.CreateUserRequest
			user    userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			request = models.CreateUserRequest{
				Email:    "tony@stark.com",
				Image:    nil,
				Lastname: "Stark",
				Name:     "Tony",
				Password: "jarvis3000",
				UserID:   "ironman",
			}

			user = userDAO.User{
				UserID:      "admin",
				BodyFat:     5,
				CurrentWeek: []bool{true, true, false, true, false, false, false},
				Email:       "admin@admin.com",
				Experience:  100,
				LastName:    "Wick",
				Name:        "John",
				Password:    "admin123",
				Streak:      10,
				Weight:      80,
			}
		})

		It("CASE: Successful user creation", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().GetUserByEmail(request.Email, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().CreateUser(gomock.Any(), ctxLogger).
				Times(1).
				Return(nil)

			mockSessionService.EXPECT().GenerateSession(request.UserID).
				Times(1).
				Return("jwt-token", nil)

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Token).To(Equal("jwt-token"))
		})

		It("CASE: User creation failed cause user already exist", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(customErrors.ConflictError{}))
		})

		It("CASE: User creation failed cause email already exist", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().GetUserByEmail(request.Email, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(customErrors.ConflictError{}))
		})

		It("CASE: User creation failed when processing get user database error", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(nil, errors.New("panic"))

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("CASE: User creation failed when processing get user by email database error", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().GetUserByEmail(request.Email, ctxLogger).
				Times(1).
				Return(nil, errors.New("panic"))

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("CASE: User creation failed when processing user creation database error", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().GetUserByEmail(request.Email, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().CreateUser(gomock.Any(), ctxLogger).
				Times(1).
				Return(errors.New("panic"))

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("CASE: User creation failed when processing a session service error", func() {

			mockUserDAO.EXPECT().GetUser(request.UserID, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().GetUserByEmail(request.Email, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			mockUserDAO.EXPECT().CreateUser(gomock.Any(), ctxLogger).
				Times(1).
				Return(nil)

			mockSessionService.EXPECT().GenerateSession(request.UserID).
				Times(1).
				Return("", errors.New("panic"))

			response, err := service.CreateUser(&request, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

	})

})
