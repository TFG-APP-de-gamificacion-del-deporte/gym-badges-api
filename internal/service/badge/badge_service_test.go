package badge_service

import (
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"
	"time"

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
				AnyTimes().
				Return(&user, nil)

			mockBadgeDAO.EXPECT().GetBadges(ctxLogger).
				Times(1).
				Return(badges, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(4), nil)

			mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
				AnyTimes().
				Return(false, nil)

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				AnyTimes().
				Return(&user, nil)

			mockBadgeDAO.EXPECT().GetBadge(gomock.Any(), ctxLogger).
				AnyTimes().
				Return(&badgeDAO.Badge{
					ID:            66,
					Name:          "First Week",
					Description:   "First week streak",
					Image:         "/image-66.jpg",
					ParentBadgeID: 6,
					Exp:           100,
				}, nil)

			mockUserDAO.EXPECT().GetUserWithGlobalRank(userID, ctxLogger).
				Times(1).
				Return(&user, int64(1000), nil)

			mockUserDAO.EXPECT().GetAttendanceCount(userID, ctxLogger).
				Times(1).
				Return(int32(0), nil)

			response, err := service.GetBadgesByUserID(userID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response)).To(Equal(2))

			Expect(response[0].ID).To(Equal(int32(1)))
			Expect(response[0].Achieved).To(Equal(true))
			Expect(response[0].Description).To(Equal("Badge 1"))
			Expect(response[0].Image).To(Equal("/image-1.jpg"))
			Expect(response[0].Name).To(Equal("badge1"))

			Expect(response[0].Children[0].ID).To(Equal(int32(2)))
			Expect(response[0].Children[0].Achieved).To(Equal(true))
			Expect(response[0].Children[0].Description).To(Equal("Badge 2"))
			Expect(response[0].Children[0].Image).To(Equal("/image-2.jpg"))
			Expect(response[0].Children[0].Name).To(Equal("badge2"))
			Expect(response[0].Children[0].Children[0].ID).To(Equal(int32(4)))
			Expect(response[0].Children[0].Children[0].Achieved).To(Equal(false))
			Expect(response[0].Children[0].Children[0].Description).To(Equal("Badge 4"))
			Expect(response[0].Children[0].Children[0].Image).To(Equal("/image-4.jpg"))
			Expect(response[0].Children[0].Children[0].Name).To(Equal("badge4"))
			Expect(response[0].Children[0].Children[0].Children).To(BeEmpty())

			Expect(response[0].Children[1].Achieved).To(Equal(false))
			Expect(response[0].Children[1].Description).To(Equal("Badge 3"))
			Expect(response[0].Children[1].Image).To(Equal("/image-3.jpg"))
			Expect(response[0].Children[1].Name).To(Equal("badge3"))
			Expect(response[0].Children[1].Children[0].ID).To(Equal(int32(5)))
			Expect(response[0].Children[1].Children[0].Description).To(Equal("Badge 5"))
			Expect(response[0].Children[1].Children[0].Image).To(Equal("/image-5.jpg"))
			Expect(response[0].Children[1].Children[0].Name).To(Equal("badge5"))
			Expect(response[0].Children[1].Children[0].Children).To(BeEmpty())

			Expect(response[1].ID).To(Equal(int32(6)))
			Expect(response[1].Achieved).To(Equal(false))
			Expect(response[1].Description).To(Equal("Badge 6"))
			Expect(response[1].Image).To(Equal("/image-6.jpg"))
			Expect(response[1].Name).To(Equal("badge6"))
			Expect(response[1].Children).To(BeEmpty())
		})

		It("CASE: Successful retrieval for user without badges", func() {

			user.Badges = nil

			mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
				AnyTimes().
				Return(&user, nil)

			mockBadgeDAO.EXPECT().GetBadges(ctxLogger).
				Times(1).
				Return(badges, nil)

			mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
				Times(1).
				Return(int32(4), nil)

			mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
				AnyTimes().
				Return(false, nil)

			mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
				AnyTimes().
				Return(&user, nil)

			mockUserDAO.EXPECT().GetUserWithGlobalRank(userID, ctxLogger).
				Times(1).
				Return(&user, int64(1000), nil)

			mockUserDAO.EXPECT().GetAttendanceCount(userID, ctxLogger).
				Times(1).
				Return(int32(0), nil)

			mockBadgeDAO.EXPECT().GetBadge(gomock.Any(), ctxLogger).
				AnyTimes().
				Return(&badgeDAO.Badge{
					ID:            66,
					Name:          "First Week",
					Description:   "First week streak",
					Image:         "/image-66.jpg",
					ParentBadgeID: 6,
					Exp:           100,
				}, nil)

			response, err := service.GetBadgesByUserID(userID, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response)).To(Equal(2))

			Expect(response[0].ID).To(Equal(int32(1)))
			Expect(response[0].Achieved).To(Equal(false))

			Expect(response[0].Children[0].ID).To(Equal(int32(2)))
			Expect(response[0].Children[0].Achieved).To(Equal(false))

			Expect(response[0].Children[0].Children[0].ID).To(Equal(int32(4)))
			Expect(response[0].Children[0].Children[0].Achieved).To(Equal(false))

			Expect(response[0].Children[1].ID).To(Equal(int32(3)))
			Expect(response[0].Children[1].Achieved).To(Equal(false))

			Expect(response[0].Children[1].Children[0].ID).To(Equal(int32(5)))
			Expect(response[0].Children[1].Children[0].Achieved).To(Equal(false))

			Expect(response[1].ID).To(Equal(int32(6)))
			Expect(response[1].Achieved).To(Equal(false))
		})

	})

	Context("Auto Badges", func() {
		var (
			ctxLogger *log.Entry
			userID    string
			user      userDAO.User
			badgeSvc  *badgesService
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()
			userID = "admin"
			user = userDAO.User{
				ID:        "admin",
				Email:     "admin@admin.com",
				Name:      "John",
				CreatedAt: time.Now().Add(-24 * time.Hour), // 1 día de antigüedad
				Streak:    0,
				Badges: []*badgeDAO.Badge{
					{
						ID:          6,
						Name:        "",
						Description: "",
						Image:       "",
						Exp:         0,
					},
				},
			}
			badgeSvc = service.(*badgesService)
		})

		Context("Streak Badges", func() {

			It("should not award any streak badges when streak is 0", func() {
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				err := badgeSvc.checkStreakBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

			It("should award first week badge when streak is 1", func() {

				user.Streak = 1

				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(66), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(66), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            66,
						Name:          "First Week",
						Description:   "First week streak",
						Image:         "/image-66.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(66), ctxLogger).
					Times(1).
					Return(nil)

				mockUserDAO.EXPECT().AddExperience(userID, int64(100), ctxLogger).
					Times(1).
					Return(nil)

				err := badgeSvc.checkStreakBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})
		})

		Context("Attendance Badges", func() {

			It("should not award any attendance badges when count is 0", func() {
				mockUserDAO.EXPECT().GetAttendanceCount(userID, ctxLogger).
					Times(1).
					Return(int32(0), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				err := badgeSvc.checkAttendancesBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

			It("should award first attendance badge when count is 100", func() {

				mockUserDAO.EXPECT().GetAttendanceCount(userID, ctxLogger).
					Times(1).
					Return(int32(100), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(72), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(72), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            72,
						Name:          "First 100",
						Description:   "First 100 gym sessions",
						Image:         "/image-72.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(72), ctxLogger).
					Times(1).
					Return(nil)

				mockUserDAO.EXPECT().AddExperience(userID, int64(100), ctxLogger).
					Times(1).
					Return(nil)

				err := badgeSvc.checkAttendancesBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})
		})

		Context("Time Badges", func() {
			It("should not award any time badges when user is new", func() {
				user.CreatedAt = time.Now()
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				err := badgeSvc.checkTimeBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

			It("should award first month badge when user is 30 days old", func() {
				user.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(71), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(71), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            71,
						Name:          "First Month",
						Description:   "First month in the gym",
						Image:         "/image-71.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(71), ctxLogger).
					Times(1).
					Return(nil)

				mockUserDAO.EXPECT().AddExperience(userID, int64(100), ctxLogger).
					Times(1).
					Return(nil)

				err := badgeSvc.checkTimeBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})
		})

		Context("Global Ranking Badges", func() {
			It("should not award any global ranking badges when rank is too low", func() {
				mockUserDAO.EXPECT().GetUserWithGlobalRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(1000), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				err := badgeSvc.checkGlobalRankingBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

			It("should award top 500 badge when rank is 500", func() {
				mockUserDAO.EXPECT().GetUserWithGlobalRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(500), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(80), ctxLogger).
					Times(1).
					Return(false, nil)

				// Mocks para AddBadge
				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(80), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            80,
						Name:          "Top 500",
						Description:   "Get to the top 500 at the Global Ranking",
						Image:         "/image-80.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(80), ctxLogger).
					Times(1).
					Return(nil)

				mockUserDAO.EXPECT().AddExperience(userID, int64(100), ctxLogger).
					Times(1).
					Return(nil)

				err := badgeSvc.checkGlobalRankingBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})
		})

		Context("Friends Ranking Badges", func() {

			It("should not award any friends ranking badges when friends count is too low", func() {
				mockUserDAO.EXPECT().GetUserWithFriendsRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(1), nil)

				mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
					Times(1).
					Return(int32(4), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				err := badgeSvc.checkFriendsRankingBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

		})

		Context("Friend Count Badges", func() {
			It("should not award any friend count badges when count is too low", func() {
				mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
					Times(1).
					Return(int32(0), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				err := badgeSvc.checkFriendCountBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

			It("should award first friend badge when count is 1", func() {
				mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
					Times(1).
					Return(int32(1), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(86), ctxLogger).
					Times(1).
					Return(false, nil)

				// Mocks para AddBadge
				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(86), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            86,
						Name:          "First Friend",
						Description:   "Add your first friend",
						Image:         "/image-86.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(86), ctxLogger).
					Times(1).
					Return(nil)

				mockUserDAO.EXPECT().AddExperience(userID, int64(100), ctxLogger).
					Times(1).
					Return(nil)

				err := badgeSvc.checkFriendCountBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})
		})

		Context("Auto Badges Coordinator", func() {
			It("should check all badge types in parallel", func() {
				// Streak Badges
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				// Attendance Badges
				mockUserDAO.EXPECT().GetAttendanceCount(userID, ctxLogger).
					Times(1).
					Return(int32(0), nil)

				// Time Badges
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				// Global Ranking Badges
				mockUserDAO.EXPECT().GetUserWithGlobalRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(1000), nil)

				// Friend Count Badges
				mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
					Times(2).
					Return(int32(0), nil)

				// Friends Ranking Badges
				mockUserDAO.EXPECT().GetUserWithFriendsRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(1000), nil)

				err := badgeSvc.checkAutoBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})

			It("should award multiple badges when conditions are met", func() {
				// Streak Badges
				user.Streak = 1
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(66), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(66), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            66,
						Name:          "First Week",
						Description:   "First week streak",
						Image:         "/image-66.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(66), ctxLogger).
					Times(1).
					Return(nil)

				mockUserDAO.EXPECT().AddExperience(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(nil)

				// Attendance Badges
				mockUserDAO.EXPECT().GetAttendanceCount(userID, ctxLogger).
					Times(1).
					Return(int32(100), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(72), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(72), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            72,
						Name:          "First 100",
						Description:   "First 100 gym sessions",
						Image:         "/image-72.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(72), ctxLogger).
					Times(1).
					Return(nil)

				// Time Badges
				user.CreatedAt = time.Now().Add(-31 * 24 * time.Hour)
				mockUserDAO.EXPECT().GetUser(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(71), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(71), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            71,
						Name:          "First Month",
						Description:   "First month in the gym",
						Image:         "/image-71.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(71), ctxLogger).
					Times(1).
					Return(nil)

				// Global Ranking Badges
				mockUserDAO.EXPECT().GetUserWithGlobalRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(500), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(80), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(80), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            80,
						Name:          "Top 500",
						Description:   "Get to the top 500 at the Global Ranking",
						Image:         "/image-80.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(80), ctxLogger).
					Times(1).
					Return(nil)

				// Friend Count Badges
				mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
					Times(1).
					Return(int32(1), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, int16(86), ctxLogger).
					Times(1).
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					Times(1).
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(int16(86), ctxLogger).
					Times(1).
					Return(&badgeDAO.Badge{
						ID:            86,
						Name:          "First Friend",
						Description:   "Add your first friend",
						Image:         "/image-86.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, int16(86), ctxLogger).
					Times(1).
					Return(nil)

				// Friends Ranking Badges
				mockUserDAO.EXPECT().GetUserWithFriendsRank(userID, ctxLogger).
					Times(1).
					Return(&user, int64(1), nil)

				mockUserDAO.EXPECT().GetFriendsCount(userID, ctxLogger).
					Times(1).
					Return(int32(20), nil)

				mockBadgeDAO.EXPECT().CheckBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(false, nil)

				mockUserDAO.EXPECT().GetUserWithBadges(userID, ctxLogger).
					AnyTimes().
					Return(&user, nil)

				mockBadgeDAO.EXPECT().GetBadge(gomock.Any(), ctxLogger).
					AnyTimes().
					Return(&badgeDAO.Badge{
						ID:            91,
						Name:          "Top Friend",
						Description:   "Get to the top at the Friends Ranking with at least 20 friends",
						Image:         "/image-91.jpg",
						ParentBadgeID: 6,
						Exp:           100,
					}, nil)

				mockBadgeDAO.EXPECT().AddBadge(userID, gomock.Any(), ctxLogger).
					AnyTimes().
					Return(nil)

				err := badgeSvc.checkAutoBadges(userID, ctxLogger)
				Expect(err).To(BeNil())
			})
		})
	})
})
