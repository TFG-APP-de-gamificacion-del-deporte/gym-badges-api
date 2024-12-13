package friends_service

import (
	"fmt"
	configs "gym-badges-api/config/gym-badges-server"
	customErrors "gym-badges-api/internal/custom-errors"
	badgeDAO "gym-badges-api/internal/repository/badge"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"gym-badges-api/tools/utils"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestServiceFriendsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: Friends Test Suite")
}

var _ = Describe("SERVICE: Friends Test Suite", func() {

	var (
		mockCtrl    *gomock.Controller
		mockUserDAO *mockDAO.MockIUserDAO
		service     IFriendsService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		service = NewFriendsService(mockUserDAO)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("Get Friends by user_id", func() {

		var (
			ctxLogger *log.Entry
			userID    string
			page      int32
			user      userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			userID = "admin"
			page = 1
			configs.Basic.FriendsPageSize = 3

			user = userDAO.User{
				ID:      "admin",
				Email:   "admin@admin.com",
				Name:    "John",
				Friends: buildFriends(5),
			}
		})

		It("CASE: Successful get friends by user_id", func() {

			offset := int32(0)

			mockUserDAO.EXPECT().GetUserWithFriends(userID, offset, configs.Basic.FriendsPageSize, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetFriendsByUserID(userID, page, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Friends)).To(Equal(5))
			Expect(response.Friends[0].Name).To(Equal("user-0"))
			Expect(response.Friends[0].Fat).To(Equal(float32(0)))
			Expect(response.Friends[0].Image).To(Equal("/image-0.jpg"))
			Expect(response.Friends[0].Streak).To(Equal(int32(0)))
			Expect(response.Friends[0].User).To(Equal("user-0"))
			Expect(response.Friends[0].Weight).To(Equal(float32(0)))
			Expect(response.Friends[0].Level).To(Equal(int32(0)))
			Expect(response.Friends[0].TopFeats[0].Image).To(Equal("/image-0.jpg"))
			Expect(response.Friends[0].TopFeats[0].Description).To(Equal("description-0"))
			Expect(response.Friends[0].TopFeats[0].Name).To(Equal("badge-0"))

			Expect(response.Friends[1].Name).To(Equal("user-1"))
			Expect(response.Friends[1].Fat).To(Equal(float32(1)))
			Expect(response.Friends[1].Image).To(Equal("/image-1.jpg"))
			Expect(response.Friends[1].Streak).To(Equal(int32(10)))
			Expect(response.Friends[1].User).To(Equal("user-1"))
			Expect(response.Friends[1].Weight).To(Equal(float32(1.5)))
			Expect(response.Friends[1].Level).To(Equal(int32(10)))
			Expect(response.Friends[1].TopFeats[0].Image).To(Equal("/image-1.jpg"))
			Expect(response.Friends[1].TopFeats[0].Description).To(Equal("description-1"))
			Expect(response.Friends[1].TopFeats[0].Name).To(Equal("badge-1"))

			Expect(response.Friends[4].Name).To(Equal("user-4"))
			Expect(response.Friends[4].Fat).To(Equal(float32(4)))
			Expect(response.Friends[4].Image).To(Equal("/image-4.jpg"))
			Expect(response.Friends[4].Streak).To(Equal(int32(40)))
			Expect(response.Friends[4].User).To(Equal("user-4"))
			Expect(response.Friends[4].Weight).To(Equal(float32(6)))
			Expect(response.Friends[4].Level).To(Equal(int32(40)))
			Expect(response.Friends[4].TopFeats[0].Image).To(Equal("/image-4.jpg"))
			Expect(response.Friends[4].TopFeats[0].Description).To(Equal("description-4"))
			Expect(response.Friends[4].TopFeats[0].Name).To(Equal("badge-4"))
		})

		It("CASE: Successful retrieval without friends", func() {

			user.Friends = nil

			page = 2
			offset := int32(3)

			mockUserDAO.EXPECT().GetUserWithFriends(userID, offset, configs.Basic.FriendsPageSize, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetFriendsByUserID(userID, page, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Friends)).To(Equal(0))
		})

		It("CASE: Get weight history failed cause user not exist", func() {

			offset := int32(0)

			mockUserDAO.EXPECT().GetUserWithFriends(userID, offset, configs.Basic.FriendsPageSize, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetFriendsByUserID(userID, page, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})

	})

})

func buildFriends(num int) []*userDAO.User {

	friends := make([]*userDAO.User, num)

	for i := 0; i < num; i++ {

		friends[i] = &userDAO.User{
			ID:         fmt.Sprintf("user-%d", i),
			BodyFat:    utils.NewFloat32(1.0 * float32(i)),
			Email:      fmt.Sprintf("user-%d@local.com", i),
			Experience: int64(i) * 1000,
			Image:      fmt.Sprintf("/image-%d.jpg", i),
			Name:       fmt.Sprintf("user-%d", i),
			Streak:     int32(i) * 10,
			Weight:     utils.NewFloat32(float32(i) * 1.5),
			TopFeats: []*badgeDAO.Badge{
				{
					Name:        fmt.Sprintf("badge-%d", i),
					Description: fmt.Sprintf("description-%d", i),
					Image:       fmt.Sprintf("/image-%d.jpg", i),
				},
			},
		}
	}

	return friends
}
