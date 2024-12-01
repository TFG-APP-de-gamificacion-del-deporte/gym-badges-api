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
				{
					ID:            4,
					Name:          "badge4",
					Description:   "Badge 4",
					Image:         "/image-4.jpg",
					ParentBadgeID: 2,
				},
				{
					ID:            5,
					Name:          "badge5",
					Description:   "Badge 5",
					Image:         "/image-5.jpg",
					ParentBadgeID: 3,
				},
				{
					ID:            6,
					Name:          "badge6",
					Description:   "Badge 6",
					Image:         "/image-6.jpg",
					ParentBadgeID: 0,
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
			Expect(len(response)).To(Equal(2))

			Expect(response[0].ID).To(Equal("1"))
			Expect(response[0].Achieved).To(Equal(true))
			Expect(response[0].Description).To(Equal("Badge 1"))
			Expect(response[0].Image).To(Equal("/image-1.jpg"))
			Expect(response[0].Name).To(Equal("badge1"))

			Expect(response[0].Children[0].ID).To(Equal("2"))
			Expect(response[0].Children[0].Achieved).To(Equal(true))
			Expect(response[0].Children[0].Description).To(Equal("Badge 2"))
			Expect(response[0].Children[0].Image).To(Equal("/image-2.jpg"))
			Expect(response[0].Children[0].Name).To(Equal("badge2"))
			Expect(response[0].Children[0].Children[0].ID).To(Equal("4"))
			Expect(response[0].Children[0].Children[0].Achieved).To(Equal(false))
			Expect(response[0].Children[0].Children[0].Description).To(Equal("Badge 4"))
			Expect(response[0].Children[0].Children[0].Image).To(Equal("/image-4.jpg"))
			Expect(response[0].Children[0].Children[0].Name).To(Equal("badge4"))
			Expect(response[0].Children[0].Children[0].Children).To(BeNil())

			Expect(response[0].Children[1].Achieved).To(Equal(false))
			Expect(response[0].Children[1].Description).To(Equal("Badge 3"))
			Expect(response[0].Children[1].Image).To(Equal("/image-3.jpg"))
			Expect(response[0].Children[1].Name).To(Equal("badge3"))
			Expect(response[0].Children[1].Children[0].ID).To(Equal("5"))
			Expect(response[0].Children[1].Children[0].Description).To(Equal("Badge 5"))
			Expect(response[0].Children[1].Children[0].Image).To(Equal("/image-5.jpg"))
			Expect(response[0].Children[1].Children[0].Name).To(Equal("badge5"))
			Expect(response[0].Children[1].Children[0].Children).To(BeNil())

			Expect(response[1].ID).To(Equal("6"))
			Expect(response[1].Achieved).To(Equal(false))
			Expect(response[1].Description).To(Equal("Badge 6"))
			Expect(response[1].Image).To(Equal("/image-6.jpg"))
			Expect(response[1].Name).To(Equal("badge6"))
			Expect(response[1].Children).To(BeNil())
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
			Expect(len(response)).To(Equal(2))

			Expect(response[0].ID).To(Equal("1"))
			Expect(response[0].Achieved).To(Equal(false))

			Expect(response[0].Children[0].ID).To(Equal("2"))
			Expect(response[0].Children[0].Achieved).To(Equal(false))

			Expect(response[0].Children[0].Children[0].ID).To(Equal("4"))
			Expect(response[0].Children[0].Children[0].Achieved).To(Equal(false))

			Expect(response[0].Children[1].ID).To(Equal("3"))
			Expect(response[0].Children[1].Achieved).To(Equal(false))

			Expect(response[0].Children[1].Children[0].ID).To(Equal("5"))
			Expect(response[0].Children[1].Children[0].Achieved).To(Equal(false))

			Expect(response[1].ID).To(Equal("6"))
			Expect(response[1].Achieved).To(Equal(false))
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
