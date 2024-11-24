package badge_service

import (
	"errors"
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestServiceBadgeSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: Badge Test Suite")
}

var _ = Describe("SERVICE: Badge Test Suite", func() {

	var (
		mockCtrl     *gomock.Controller
		mockUserDAO  *mockDAO.MockIUserDAO
		mockBadgeDAO *mockDAO.MockIBadgeDAO
		service      IBadgeService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		mockBadgeDAO = mockDAO.NewMockIBadgeDAO(mockCtrl)
		service = NewBadgeService(mockUserDAO, mockBadgeDAO)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("Get Badge by user_id", func() {

		var (
			ctxLogger *log.Entry
			userID    string
			user      userDAO.User
			badges    []*badgeDAO.Badge
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			userID = "admin"

			badges = []*badgeDAO.Badge{
				{
					ID:            1,
					Name:          "badge1",
					Description:   "Badge 1",
					Image:         "/image-1.jpg",
					ParentBadgeID: 0,
					ParentBadge:   nil,
				},
				{
					ID:            2,
					Name:          "badge2",
					Description:   "Badge 2",
					Image:         "/image-2.jpg",
					ParentBadgeID: 1,
				},
				{
					ID:            3,
					Name:          "badge3",
					Description:   "Badge 3",
					Image:         "/image-3.jpg",
					ParentBadgeID: 1,
				},
			}

			user = userDAO.User{
				ID:     "admin",
				Email:  "admin@admin.com",
				Name:   "John",
				Badges: badges[:2],
			}
		})

		It("CASE: Successful get badges by user_id", func() {

			mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockBadgeDAO.EXPECT().GetBadges(ctxLogger).
				Times(1).
				Return(badges, nil)

			response, err := service.GetBadgesByUserID(userID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Badges)).To(Equal(3))

			Expect(response.Badges["1"].Achieved).To(Equal(true))
			Expect(response.Badges["1"].Description).To(Equal("Badge 1"))
			Expect(response.Badges["1"].Image).To(Equal("/image-1.jpg"))
			Expect(response.Badges["1"].Name).To(Equal("badge1"))
			Expect(response.Badges["1"].Parent).To(Equal("0"))

			Expect(response.Badges["2"].Achieved).To(Equal(true))
			Expect(response.Badges["2"].Description).To(Equal("Badge 2"))
			Expect(response.Badges["2"].Image).To(Equal("/image-2.jpg"))
			Expect(response.Badges["2"].Name).To(Equal("badge2"))
			Expect(response.Badges["2"].Parent).To(Equal("1"))

			Expect(response.Badges["3"].Achieved).To(Equal(false))
			Expect(response.Badges["3"].Description).To(Equal("Badge 3"))
			Expect(response.Badges["3"].Image).To(Equal("/image-3.jpg"))
			Expect(response.Badges["3"].Name).To(Equal("badge3"))
			Expect(response.Badges["3"].Parent).To(Equal("1"))
		})

		It("CASE: Successful retrieval for user without badges", func() {

			user.Badges = nil

			mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockBadgeDAO.EXPECT().GetBadges(ctxLogger).
				Times(1).
				Return(badges, nil)

			response, err := service.GetBadgesByUserID(userID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Badges)).To(Equal(3))

			Expect(response.Badges["1"].Achieved).To(Equal(false))
			Expect(response.Badges["2"].Achieved).To(Equal(false))
			Expect(response.Badges["3"].Achieved).To(Equal(false))
		})

		It("CASE: Get user badges failed cause user dao respond with a error", func() {

			mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
				Times(1).
				Return(nil, errors.New("timeout"))

			mockBadgeDAO.EXPECT().GetBadges(ctxLogger).
				Times(1).
				Return(badges, nil)

			response, err := service.GetBadgesByUserID(userID, ctxLogger)
			Expect(err).To(Not(BeNil()))
			Expect(response).To(BeNil())
		})

		It("CASE: Get user badges failed cause badges dao respond with a error", func() {

			mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
				Times(1).
				Return(&user, nil)

			mockBadgeDAO.EXPECT().GetBadges(ctxLogger).
				Times(1).
				Return(nil, errors.New("timeout"))

			response, err := service.GetBadgesByUserID(userID, ctxLogger)
			Expect(err).To(Not(BeNil()))
			Expect(response).To(BeNil())
		})

	})

})
