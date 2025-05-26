package user_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	mockService "gym-badges-api/mocks/service"
	"gym-badges-api/models"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"gym-badges-api/tools/utils"
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
		mockBadgeDAO       *mockDAO.MockIBadgeDAO
		mockSessionService *mockService.MockISessionService
		service            IUserService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		mockBadgeDAO = mockDAO.NewMockIBadgeDAO(mockCtrl)
		mockSessionService = mockService.NewMockISessionService(mockCtrl)
		service = NewUserService(mockUserDAO, mockBadgeDAO, mockSessionService)
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
				Name:     "Tony",
				Password: "jarvis3000",
				UserID:   "ironman",
			}

			user = userDAO.User{
				ID:          "admin",
				BodyFat:     utils.NewFloat32(5),
				CurrentWeek: []bool{true, true, false, true, false, false, false},
				Email:       "admin@admin.com",
				Experience:  100,
				Name:        "John",
				Password:    "admin123",
				Streak:      10,
				Weight:      utils.NewFloat32(80),
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

	Context("Get User", func() {
		var (
			ctxLogger  *log.Entry
			userID     string
			authUserID string
			user       userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()
			userID = "admin"
			authUserID = "admin"
			user = userDAO.User{
				ID:          "admin",
				BodyFat:     utils.NewFloat32(5),
				CurrentWeek: []bool{true, true, false, true, false, false, false},
				Email:       "admin@admin.com",
				Experience:  100,
				Name:        "John",
				Password:    "admin123",
				Streak:      10,
				Weight:      utils.NewFloat32(80),
				Height:      utils.NewFloat32(180),
				Sex:         "M",
				WeeklyGoal:  3,
				Preferences: []userDAO.Preference{
					{ID: 1, On: false, UserID: "admin"},
					{ID: 2, On: false, UserID: "admin"},
				},
				TopFeats: []*badgeDAO.Badge{
					{ID: 1, Name: "Feat 1", Description: "Description 1", Image: "/image1.jpg"},
					{ID: 2, Name: "Feat 2", Description: "Description 2", Image: "/image2.jpg"},
				},
			}
		})

		It("CASE: Successful get user info for own profile", func() {
			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(5), nil)

			response, err := service.GetUser(userID, authUserID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.UserID).To(Equal(user.ID))
			Expect(response.Name).To(Equal(user.Name))
			Expect(response.Experience).To(Equal(user.Experience))
			Expect(response.Streak).To(Equal(user.Streak))
			Expect(response.IsFriend).To(BeTrue())
			Expect(response.TotalFriends).To(Equal(int32(5)))
			Expect(len(response.Preferences)).To(Equal(2))
			Expect(len(response.TopFeats)).To(Equal(2))
		})

		It("CASE: Successful get user info for other profile with private account", func() {
			authUserID = "other"
			user.Preferences[0].On = true // Private account

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(5), nil)

			mockUserDAO.EXPECT().CheckFriendship(userID, authUserID, ctxLogger).
				Times(1).
				Return(false, nil)

			response, err := service.GetUser(userID, authUserID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.UserID).To(Equal(user.ID))
			Expect(response.Name).To(Equal(user.Name))
			Expect(response.Experience).To(Equal(int64(0)))
			Expect(response.Streak).To(Equal(int32(0)))
			Expect(response.IsFriend).To(BeFalse())
			Expect(response.TotalFriends).To(Equal(int32(5)))
		})

		It("CASE: Successful get user info for other profile with hidden weight", func() {
			authUserID = "other"
			user.Preferences[1].On = true // Hide weight and fat

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(5), nil)

			mockUserDAO.EXPECT().CheckFriendship(userID, authUserID, ctxLogger).
				Times(1).
				Return(true, nil)

			response, err := service.GetUser(userID, authUserID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.UserID).To(Equal(user.ID))
			Expect(response.Name).To(Equal(user.Name))
			Expect(response.Experience).To(Equal(user.Experience))
			Expect(response.Streak).To(Equal(user.Streak))
			Expect(response.IsFriend).To(BeTrue())
			Expect(response.TotalFriends).To(Equal(int32(5)))
			Expect(response.BodyFat).To(BeNil())
			Expect(response.Weight).To(BeNil())
		})

		It("CASE: Get user info failed when getting user", func() {
			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(nil, errors.New("database error"))

			response, err := service.GetUser(userID, authUserID, ctxLogger)
			Expect(err).ToNot(BeNil())
			Expect(response).To(BeNil())
		})

		It("CASE: Get user info failed when getting friends count", func() {
			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(0), errors.New("database error"))

			response, err := service.GetUser(userID, authUserID, ctxLogger)
			Expect(err).ToNot(BeNil())
			Expect(response).To(BeNil())
		})

		It("CASE: Get user info failed when checking friendship", func() {
			authUserID = "other"

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(5), nil)

			mockUserDAO.EXPECT().CheckFriendship(userID, authUserID, ctxLogger).
				Times(1).
				Return(false, errors.New("database error"))

			response, err := service.GetUser(userID, authUserID, ctxLogger)
			Expect(err).ToNot(BeNil())
			Expect(response).To(BeNil())
		})
	})

	Context("Edit User Info", func() {
		var (
			ctxLogger *log.Entry
			userID    string
			request   models.EditUserInfoRequest
			user      userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()
			userID = "admin"
			request = models.EditUserInfoRequest{
				Email:      "new@email.com",
				Name:       "New Name",
				WeeklyGoal: 4,
				Height:     175,
				Sex:        "F",
				TopFeats:   []int32{1, 2, 3},
				Preferences: []*models.Preference{
					{PreferenceID: 1, On: true},
					{PreferenceID: 2, On: false},
				},
			}
			user = userDAO.User{
				ID:         "admin",
				Email:      "new@email.com",
				Name:       "New Name",
				WeeklyGoal: 4,
				Height:     utils.NewFloat32(175),
				Sex:        "F",
				TopFeats: []*badgeDAO.Badge{
					{ID: 1},
					{ID: 2},
					{ID: 3},
				},
				Preferences: []userDAO.Preference{
					{ID: 1, On: true, UserID: "admin"},
					{ID: 2, On: false, UserID: "admin"},
				},
			}
		})

		It("CASE: Successful edit user info", func() {
			mockUserDAO.EXPECT().EditUserInfo(userID, gomock.Any(), ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.EditUserInfo(userID, &request, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.UserID).To(Equal(user.ID))
			Expect(response.Name).To(Equal(user.Name))
			Expect(response.Image).To(Equal(user.Image))
			Expect(response.Height).To(Equal(*user.Height))
			Expect(response.Sex).To(Equal(user.Sex))
			Expect(len(response.TopFeats)).To(Equal(3))
			Expect(len(response.Preferences)).To(Equal(2))
		})

		It("CASE: Edit user info failed", func() {
			mockUserDAO.EXPECT().EditUserInfo(userID, gomock.Any(), ctxLogger).
				Times(1).
				Return(nil, errors.New("database error"))

			response, err := service.EditUserInfo(userID, &request, ctxLogger)
			Expect(err).ToNot(BeNil())
			Expect(response).To(BeNil())
		})
	})

	Context("Edit User Preferences", func() {
		var (
			ctxLogger *log.Entry
			userID    string
			request   models.EditUserPreferenceRequest
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()
			userID = "admin"
			request = models.EditUserPreferenceRequest{
				PrivateAccount:   true,
				HideWeightAndFat: false,
			}
		})

		It("CASE: Successful edit user preferences", func() {
			mockUserDAO.EXPECT().UpdateUserPreferences(userID, gomock.Any(), ctxLogger).
				Times(1).
				Return(nil)

			err := service.EditUserPreferences(userID, &request, ctxLogger)
			Expect(err).To(BeNil())
		})

		It("CASE: Edit user preferences failed", func() {
			mockUserDAO.EXPECT().UpdateUserPreferences(userID, gomock.Any(), ctxLogger).
				Times(1).
				Return(errors.New("database error"))

			err := service.EditUserPreferences(userID, &request, ctxLogger)
			Expect(err).ToNot(BeNil())
		})
	})

	Context("Edit User Top Feats", func() {
		var (
			ctxLogger *log.Entry
			userID    string
			topFeats  []int16
			badges    []*badgeDAO.Badge
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()
			userID = "admin"
			topFeats = []int16{1, 2, 3}
			badges = []*badgeDAO.Badge{
				{ID: 1, Name: "Feat 1", Description: "Description 1", Image: "/image1.jpg"},
				{ID: 2, Name: "Feat 2", Description: "Description 2", Image: "/image2.jpg"},
				{ID: 3, Name: "Feat 3", Description: "Description 3", Image: "/image3.jpg"},
			}
		})

		It("CASE: Successful edit user top feats", func() {

			mockBadgeDAO.EXPECT().GetBadgesByIds(topFeats, ctxLogger).
				Times(1).
				Return(badges, nil)

			mockUserDAO.EXPECT().UpdateUserTopFeats(userID, badges, ctxLogger).
				Times(1).
				Return(nil)

			err := service.EditUserTopFeats(userID, topFeats, ctxLogger)
			Expect(err).To(BeNil())
		})

		It("CASE: Edit user top feats failed when getting badges", func() {
			mockBadgeDAO.EXPECT().GetBadgesByIds(topFeats, ctxLogger).
				Times(1).
				Return(nil, errors.New("database error"))

			err := service.EditUserTopFeats(userID, topFeats, ctxLogger)
			Expect(err).ToNot(BeNil())
		})

		It("CASE: Edit user top feats failed when invalid badges", func() {
			mockBadgeDAO.EXPECT().GetBadgesByIds(topFeats, ctxLogger).
				Times(1).
				Return(badges[:2], nil)

			err := service.EditUserTopFeats(userID, topFeats, ctxLogger)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("invalid top feats"))
		})

		It("CASE: Edit user top feats failed when updating", func() {
			mockBadgeDAO.EXPECT().GetBadgesByIds(topFeats, ctxLogger).
				Times(1).
				Return(badges, nil)

			mockUserDAO.EXPECT().UpdateUserTopFeats(userID, badges, ctxLogger).
				Times(1).
				Return(errors.New("database error"))

			err := service.EditUserTopFeats(userID, topFeats, ctxLogger)
			Expect(err).ToNot(BeNil())
		})
	})
})
