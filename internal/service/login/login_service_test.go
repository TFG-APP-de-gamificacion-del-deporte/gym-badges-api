package login_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	mockService "gym-badges-api/mocks/service"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestServiceLoginSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: Login Test Suite")
}

var _ = Describe("SERVICE: Login Test Suite", func() {

	var (
		mockCtrl           *gomock.Controller
		mockUserDAO        *mockDAO.MockIUserDAO
		mockSessionService *mockService.MockISessionService
		service            ILoginService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		mockSessionService = mockService.NewMockISessionService(mockCtrl)
		service = NewLoginService(mockUserDAO, mockSessionService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("Login", func() {

		var (
			ctxLogger *log.Entry

			userID   string
			password string
			user     userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			userID = "admin"
			password = "admin123"

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

		It("CASE: Successful login with valid credentials", func() {

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockSessionService.EXPECT().GenerateSession(userID).
				Times(1).
				Return("jwt-token", nil)

			response, err := service.Login(userID, password, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Token).To(Equal("jwt-token"))
		})

		It("CASE: Login failed with invalid credentials", func() {

			password = "invalid"

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.Login(userID, password, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).To(BeAssignableToTypeOf(customErrors.UnauthorizedError{}))
		})

		It("CASE: Login failed when processing a database error", func() {

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(nil, errors.New("panic"))

			response, err := service.Login(userID, password, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("CASE: Login failed when processing a session service error", func() {

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockSessionService.EXPECT().GenerateSession(userID).
				Times(1).
				Return("", errors.New("panic"))

			response, err := service.Login(userID, password, ctxLogger)
			Expect(response).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

	})

})
