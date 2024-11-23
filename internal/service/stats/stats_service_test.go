package stats_service

import (
	"fmt"
	customErrors "gym-badges-api/internal/custom-errors"
	userDAO "gym-badges-api/internal/repository/user"
	mockDAO "gym-badges-api/mocks/dao"
	mockService "gym-badges-api/mocks/service"
	toolsLogging "gym-badges-api/tools/logging"
	toolsTesting "gym-badges-api/tools/testing"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestServiceStatsSuite(t *testing.T) {
	toolsTesting.ConfigureTestSuite(t, "SERVICE: Stats Test Suite")
}

var _ = Describe("SERVICE: Stats Test Suite", func() {

	var (
		mockCtrl           *gomock.Controller
		mockUserDAO        *mockDAO.MockIUserDAO
		mockSessionService *mockService.MockISessionService
		service            IStatsService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		mockSessionService = mockService.NewMockISessionService(mockCtrl)
		service = NewStatsService(mockUserDAO, mockSessionService)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("Get Weight History", func() {

		var (
			ctxLogger *log.Entry
			userID    string
			months    int32
			user      userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			userID = "admin"
			months = 3

			user = userDAO.User{
				ID:    "admin",
				Email: "admin@admin.com",
				Name:  "John",
				WeightHistory: []userDAO.WeightHistory{
					{
						UserID: "admin",
						Date:   parseTime("2024-11-01T10:30:00"),
						Weight: 79.0,
					},
					{
						UserID: "admin",
						Date:   parseTime("2024-11-07T10:30:00"),
						Weight: 80.5,
					},
					{
						UserID: "admin",
						Date:   parseTime("2024-11-14T10:30:00"),
						Weight: 83.0,
					},
				},
			}
		})

		It("CASE: Successful get weight history", func() {

			mockUserDAO.EXPECT().GetUserWithWeightHistory(userID, months, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetWeightHistory(userID, months, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Days)).To(Equal(3))
			Expect(response.Days[0].Date).To(Equal("2024-11-01"))
			Expect(response.Days[0].Value).To(Equal(float32(79)))
			Expect(response.Days[1].Date).To(Equal("2024-11-07"))
			Expect(response.Days[1].Value).To(Equal(float32(80.5)))
			Expect(response.Days[2].Date).To(Equal("2024-11-14"))
			Expect(response.Days[2].Value).To(Equal(float32(83)))
		})

		It("CASE: Successful retrieval without weight history info", func() {

			user.WeightHistory = nil

			mockUserDAO.EXPECT().GetUserWithWeightHistory(userID, months, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetWeightHistory(userID, months, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Days)).To(Equal(0))
		})

		It("CASE: Get weight history failed cause user not exist", func() {

			user.WeightHistory = nil

			mockUserDAO.EXPECT().GetUserWithWeightHistory(userID, months, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetWeightHistory(userID, months, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})

	})

	Context("Get Fat History", func() {

		var (
			ctxLogger *log.Entry
			userID    string
			months    int32
			user      userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			userID = "admin"
			months = 3

			user = userDAO.User{
				ID:    "admin",
				Email: "admin@admin.com",
				Name:  "John",
				FatHistory: []userDAO.FatHistory{
					{
						UserID: "admin",
						Date:   parseTime("2024-11-01T10:30:00"),
						Fat:    79.0,
					},
					{
						UserID: "admin",
						Date:   parseTime("2024-11-07T10:30:00"),
						Fat:    80.5,
					},
					{
						UserID: "admin",
						Date:   parseTime("2024-11-14T10:30:00"),
						Fat:    83.0,
					},
				},
			}
		})

		It("CASE: Successful get fat history", func() {

			mockUserDAO.EXPECT().GetUserWithFatHistory(userID, months, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetFatHistory(userID, months, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Days)).To(Equal(3))
			Expect(response.Days[0].Date).To(Equal("2024-11-01"))
			Expect(response.Days[0].Value).To(Equal(float32(79)))
			Expect(response.Days[1].Date).To(Equal("2024-11-07"))
			Expect(response.Days[1].Value).To(Equal(float32(80.5)))
			Expect(response.Days[2].Date).To(Equal("2024-11-14"))
			Expect(response.Days[2].Value).To(Equal(float32(83)))
		})

		It("CASE: Successful retrieval without fat history info", func() {

			user.FatHistory = nil

			mockUserDAO.EXPECT().GetUserWithFatHistory(userID, months, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetFatHistory(userID, months, ctxLogger)
			Expect(err).To(BeNil())
			Expect(len(response.Days)).To(Equal(0))
		})

		It("CASE: Get fat history failed cause user not exist", func() {

			user.WeightHistory = nil

			mockUserDAO.EXPECT().GetUserWithFatHistory(userID, months, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetFatHistory(userID, months, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})

	})

	Context("Get Streak Calendar By Year And Month", func() {

		var (
			ctxLogger *log.Entry
			userID    string
			month     int32
			year      int32
			user      userDAO.User
		)

		BeforeEach(func() {
			ctxLogger = toolsLogging.BuildLogger()

			userID = "admin"
			year = 2024
			month = 11

			user = userDAO.User{
				ID:         "admin",
				Email:      "admin@admin.com",
				Name:       "John",
				Streak:     77,
				WeeklyGoal: 3,
				GymAttendance: []userDAO.GymAttendance{
					{
						UserID: "admin",
						Date:   parseTime("2024-11-01T10:30:00"),
					},
					{
						UserID: "admin",
						Date:   parseTime("2024-11-07T10:30:00"),
					},
					{
						UserID: "admin",
						Date:   parseTime("2024-11-14T10:30:00"),
					},
				},
			}
		})

		It("CASE: Successful get streak calendar info", func() {

			mockUserDAO.EXPECT().GetUserWithAttendance(userID, year, month, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetStreakCalendarByYearAndMonth(userID, year, month, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Streak).To(Equal(int32(77)))
			Expect(response.WeeklyGoal).To(Equal(int32(3)))
			Expect(len(response.Days)).To(Equal(3))
			Expect(response.Days[0]).To(Equal("2024-11-01"))
			Expect(response.Days[1]).To(Equal("2024-11-07"))
			Expect(response.Days[2]).To(Equal("2024-11-14"))
		})

		It("CASE: Successful retrieval without attendance info", func() {

			user.GymAttendance = nil

			mockUserDAO.EXPECT().GetUserWithAttendance(userID, year, month, ctxLogger).
				Times(1).
				Return(&user, nil)

			response, err := service.GetStreakCalendarByYearAndMonth(userID, year, month, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Streak).To(Equal(int32(77)))
			Expect(response.WeeklyGoal).To(Equal(int32(3)))
			Expect(len(response.Days)).To(Equal(0))
		})

		It("CASE: Get streak calendar failed cause user not exist", func() {

			user.WeightHistory = nil

			mockUserDAO.EXPECT().GetUserWithAttendance(userID, year, month, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetStreakCalendarByYearAndMonth(userID, year, month, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})

	})

})

func parseTime(dateStr string) time.Time {
	parsedTime, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		Fail(fmt.Sprintf("Failed to parse date: %s", dateStr), 1)
	}
	return parsedTime
}
