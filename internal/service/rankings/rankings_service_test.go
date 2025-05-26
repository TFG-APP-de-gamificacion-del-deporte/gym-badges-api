package rankings_service

import (
	"errors"
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	"gym-badges-api/models"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestServiceRankingsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: Rankings Test Suite")
}

var _ = Describe("SERVICE: Rankings Test Suite", func() {

	var (
		mockCtrl    *gomock.Controller
		mockUserDAO *mockDAO.MockIUserDAO
		service     IRankingsService
		ctxLog      *log.Entry
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		service = NewRankingsService(mockUserDAO)
		ctxLog = toolsLogging.BuildLogger()
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("GetGlobalRanking", func() {
		var (
			userID string
			page   int32
		)

		BeforeEach(func() {
			userID = "test_user"
			page = 1
		})

		It("CASE: Success Response with user in ranking", func() {
			users := []*userDAO.User{
				{
					ID:         "user1",
					Name:       "John Doe",
					Image:      []byte("image1"),
					Experience: 1000,
					Streak:     5,
				},
				{
					ID:         "user2",
					Name:       "Jane Doe",
					Image:      []byte("image2"),
					Experience: 800,
					Streak:     3,
				},
			}

			user := &userDAO.User{
				ID:         "test_user",
				Name:       "Test User",
				Image:      []byte("image3"),
				Experience: 500,
				Streak:     2,
			}

			mockUserDAO.EXPECT().GetUsersOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(users, nil)

			mockUserDAO.EXPECT().GetUserWithGlobalRank(gomock.Any(), gomock.Any()).
				Times(1).
				Return(user, int64(3), nil)

			response, err := service.GetGlobalRanking(userID, page, ctxLog)
			Expect(err).To(BeNil())
			Expect(response).To(BeEquivalentTo(&models.GetRankingResponse{
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
					Image:  []byte("image3"),
					Level:  5,
					Name:   "Test User",
					Rank:   3,
					Streak: 2,
					UserID: "test_user",
				},
			}))
		})

		It("CASE: Success Response with user not in ranking", func() {
			users := []*userDAO.User{
				{
					ID:         "user1",
					Name:       "John Doe",
					Image:      []byte("image1"),
					Experience: 1000,
					Streak:     5,
				},
				{
					ID:         "user2",
					Name:       "Jane Doe",
					Image:      []byte("image2"),
					Experience: 800,
					Streak:     3,
				},
			}

			user := &userDAO.User{
				ID:         "test_user",
				Name:       "Test User",
				Image:      []byte("image3"),
				Experience: 500,
				Streak:     2,
			}

			mockUserDAO.EXPECT().GetUsersOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(users, nil)

			mockUserDAO.EXPECT().GetUserWithGlobalRank(gomock.Any(), gomock.Any()).
				Times(1).
				Return(user, int64(15), nil)

			response, err := service.GetGlobalRanking(userID, page, ctxLog)
			Expect(err).To(BeNil())
			Expect(response).To(BeEquivalentTo(&models.GetRankingResponse{
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
			}))
		})

		It("CASE: Not Found Error Response", func() {
			mockUserDAO.EXPECT().GetUsersOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, nil)

			mockUserDAO.EXPECT().GetUserWithGlobalRank(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, int64(-1), customErrors.BuildNotFoundError("user not found"))

			response, err := service.GetGlobalRanking(userID, page, ctxLog)
			Expect(err).To(BeEquivalentTo(customErrors.BuildNotFoundError("user not found")))
			Expect(response).To(BeNil())
		})

		It("CASE: Internal Server Error Response", func() {
			mockUserDAO.EXPECT().GetUsersOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, errors.New("database error"))

			response, err := service.GetGlobalRanking(userID, page, ctxLog)
			Expect(err).To(BeEquivalentTo(errors.New("database error")))
			Expect(response).To(BeNil())
		})
	})

	Context("GetFriendsRanking", func() {
		var (
			userID string
			page   int32
		)

		BeforeEach(func() {
			userID = "test_user"
			page = 1
		})

		It("CASE: Success Response with user in ranking", func() {
			users := []*userDAO.User{
				{
					ID:         "friend1",
					Name:       "Friend One",
					Image:      []byte("image1"),
					Experience: 1000,
					Streak:     5,
				},
				{
					ID:         "friend2",
					Name:       "Friend Two",
					Image:      []byte("image2"),
					Experience: 800,
					Streak:     3,
				},
			}

			user := &userDAO.User{
				ID:         "test_user",
				Name:       "Test User",
				Image:      []byte("image3"),
				Experience: 500,
				Streak:     2,
			}

			mockUserDAO.EXPECT().GetFriendsOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(users, nil)

			mockUserDAO.EXPECT().GetUserWithFriendsRank(gomock.Any(), gomock.Any()).
				Times(1).
				Return(user, int64(3), nil)

			response, err := service.GetFriendsRanking(userID, page, ctxLog)
			Expect(err).To(BeNil())
			Expect(response).To(BeEquivalentTo(&models.GetRankingResponse{
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
					Image:  []byte("image3"),
					Level:  5,
					Name:   "Test User",
					Rank:   3,
					Streak: 2,
					UserID: "test_user",
				},
			}))
		})

		It("CASE: Success Response with user not in ranking", func() {
			users := []*userDAO.User{
				{
					ID:         "friend1",
					Name:       "Friend One",
					Image:      []byte("image1"),
					Experience: 1000,
					Streak:     5,
				},
				{
					ID:         "friend2",
					Name:       "Friend Two",
					Image:      []byte("image2"),
					Experience: 800,
					Streak:     3,
				},
			}

			user := &userDAO.User{
				ID:         "test_user",
				Name:       "Test User",
				Image:      []byte("image3"),
				Experience: 500,
				Streak:     2,
			}

			mockUserDAO.EXPECT().GetFriendsOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(users, nil)

			mockUserDAO.EXPECT().GetUserWithFriendsRank(gomock.Any(), gomock.Any()).
				Times(1).
				Return(user, int64(15), nil)

			response, err := service.GetFriendsRanking(userID, page, ctxLog)
			Expect(err).To(BeNil())
			Expect(response).To(BeEquivalentTo(&models.GetRankingResponse{
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
			}))
		})

		It("CASE: Not Found Error Response", func() {
			mockUserDAO.EXPECT().GetFriendsOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, nil)

			mockUserDAO.EXPECT().GetUserWithFriendsRank(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, int64(-1), customErrors.BuildNotFoundError("user not found"))

			response, err := service.GetFriendsRanking(userID, page, ctxLog)
			Expect(err).To(BeEquivalentTo(customErrors.BuildNotFoundError("user not found")))
			Expect(response).To(BeNil())
		})

		It("CASE: Internal Server Error Response", func() {
			mockUserDAO.EXPECT().GetFriendsOrderedByExp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil, errors.New("database error"))

			response, err := service.GetFriendsRanking(userID, page, ctxLog)
			Expect(err).To(BeEquivalentTo(errors.New("database error")))
			Expect(response).To(BeNil())
		})
	})
})
