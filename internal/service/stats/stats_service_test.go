package stats_service

import (
	"errors"
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
		ctxLogger          *log.Entry
		userID             string
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockUserDAO = mockDAO.NewMockIUserDAO(mockCtrl)
		mockSessionService = mockService.NewMockISessionService(mockCtrl)
		service = NewStatsService(mockUserDAO, mockSessionService)
		ctxLogger = toolsLogging.BuildLogger()
		userID = "test-user"
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	Context("GetWeightHistory", func() {
		It("CASE: Successfully gets weight history", func() {
			months := int32(3)
			user := &userDAO.User{
				ID: userID,
				WeightHistory: []userDAO.WeightHistory{
					{
						UserID: userID,
						Date:   time.Now(),
						Weight: 75.5,
					},
				},
			}

			mockUserDAO.EXPECT().GetUserWithWeightHistory(userID, months, ctxLogger).
				Times(1).
				Return(user, nil)

			response, err := service.GetWeightHistory(userID, months, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Days).To(HaveLen(1))
			Expect(response.Days[0].Value).To(Equal(float32(75.5)))
		})

		It("CASE: Returns error when user not found", func() {
			months := int32(3)
			mockUserDAO.EXPECT().GetUserWithWeightHistory(userID, months, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetWeightHistory(userID, months, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})
	})

	Context("AddWeight", func() {
		It("CASE: Successfully adds weight", func() {
			weight := float32(75.5)
			now := time.Now().Truncate(time.Hour * 24)

			mockUserDAO.EXPECT().AddWeight(userID, weight, now, ctxLogger).
				Times(1).
				Return(nil)

			err := service.AddWeight(userID, weight, ctxLogger)
			Expect(err).To(BeNil())
		})

		It("CASE: Returns error when adding weight fails", func() {
			weight := float32(75.5)
			now := time.Now().Truncate(time.Hour * 24)

			mockUserDAO.EXPECT().AddWeight(userID, weight, now, ctxLogger).
				Times(1).
				Return(errors.New("database error"))

			err := service.AddWeight(userID, weight, ctxLogger)
			Expect(err).ToNot(BeNil())
		})
	})

	Context("GetFatHistory", func() {
		It("CASE: Successfully gets fat history", func() {
			months := int32(3)
			user := &userDAO.User{
				ID: userID,
				FatHistory: []userDAO.FatHistory{
					{
						UserID: userID,
						Date:   time.Now(),
						Fat:    15.0,
					},
				},
			}

			mockUserDAO.EXPECT().GetUserWithFatHistory(userID, months, ctxLogger).
				Times(1).
				Return(user, nil)

			response, err := service.GetFatHistory(userID, months, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Days).To(HaveLen(1))
			Expect(response.Days[0].Value).To(Equal(float32(15.0)))
		})

		It("CASE: Returns error when user not found", func() {
			months := int32(3)
			mockUserDAO.EXPECT().GetUserWithFatHistory(userID, months, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetFatHistory(userID, months, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})
	})

	Context("AddBodyFat", func() {
		It("CASE: Successfully adds body fat", func() {
			bodyFat := float32(15.0)
			now := time.Now().Truncate(time.Hour * 24)

			mockUserDAO.EXPECT().AddBodyFat(userID, bodyFat, now, ctxLogger).
				Times(1).
				Return(nil)

			err := service.AddBodyFat(userID, bodyFat, ctxLogger)
			Expect(err).To(BeNil())
		})

		It("CASE: Returns error when adding body fat fails", func() {
			bodyFat := float32(15.0)
			now := time.Now().Truncate(time.Hour * 24)

			mockUserDAO.EXPECT().AddBodyFat(userID, bodyFat, now, ctxLogger).
				Times(1).
				Return(errors.New("database error"))

			err := service.AddBodyFat(userID, bodyFat, ctxLogger)
			Expect(err).ToNot(BeNil())
		})
	})

	Context("GetStreakCalendarByYearAndMonth", func() {
		It("CASE: Successfully gets streak calendar", func() {
			year := int32(2024)
			month := int32(3)
			user := &userDAO.User{
				ID: userID,
				GymAttendance: []userDAO.GymAttendance{
					{
						UserID: userID,
						Date:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				Streak:     5,
				WeeklyGoal: 3,
			}

			mockUserDAO.EXPECT().GetUserWithAttendance(userID, year, month, ctxLogger).
				Times(1).
				Return(user, nil)

			response, err := service.GetStreakCalendarByYearAndMonth(userID, year, month, ctxLogger)
			Expect(err).To(BeNil())
			Expect(response.Days).To(HaveLen(1))
			Expect(response.Days[0]).To(Equal("2024-03-01"))
			Expect(response.Streak).To(Equal(int32(5)))
			Expect(response.WeeklyGoal).To(Equal(int32(3)))
		})

		It("CASE: Returns error when user not found", func() {
			year := int32(2024)
			month := int32(3)
			mockUserDAO.EXPECT().GetUserWithAttendance(userID, year, month, ctxLogger).
				Times(1).
				Return(nil, customErrors.BuildNotFoundError("not found"))

			response, err := service.GetStreakCalendarByYearAndMonth(userID, year, month, ctxLogger)
			Expect(err).To(BeAssignableToTypeOf(customErrors.NotFoundError{}))
			Expect(response).To(BeNil())
		})
	})

	Context("AddGymAttendance", func() {

		It("CASE: Successfully adds gym attendance", func() {
			date := time.Now().Truncate(time.Hour * 24)
			monday := monday()
			dateIndex := int(date.Sub(monday).Hours() / 24)

			mockUserDAO.EXPECT().AddDayToCurrentWeek(userID, dateIndex, ctxLogger).
				Times(1).
				Return(nil)

			mockUserDAO.EXPECT().AddGymAttendance(userID, date, ctxLogger).
				Times(1).
				Return(nil)

			err := service.AddGymAttendance(userID, date, ctxLogger)
			Expect(err).To(BeNil())
		})

		It("CASE: Returns error when adding gym attendance fails", func() {
			date := time.Now().Truncate(time.Hour * 24)
			monday := monday()
			dateIndex := int(date.Sub(monday).Hours() / 24)

			mockUserDAO.EXPECT().AddDayToCurrentWeek(userID, dateIndex, ctxLogger).
				Times(1).
				Return(errors.New("database error"))

			err := service.AddGymAttendance(userID, date, ctxLogger)
			Expect(err).ToNot(BeNil())
		})
	})

	Context("DeleteGymAttendance", func() {
		It("CASE: Successfully deletes gym attendance", func() {
			date := time.Now().Truncate(time.Hour * 24)
			monday := monday()
			dateIndex := int(date.Sub(monday).Hours() / 24)

			mockUserDAO.EXPECT().DeleteDayFromCurrentWeek(userID, dateIndex, ctxLogger).
				Times(1).
				Return(nil)

			mockUserDAO.EXPECT().DeleteGymAttendance(userID, date, ctxLogger).
				Times(1).
				Return(nil)

			err := service.DeleteGymAttendance(userID, date, ctxLogger)
			Expect(err).To(BeNil())
		})

		It("CASE: Returns error when deleting gym attendance fails", func() {
			date := time.Now().Truncate(time.Hour * 24)
			monday := monday()
			dateIndex := int(date.Sub(monday).Hours() / 24)

			mockUserDAO.EXPECT().DeleteDayFromCurrentWeek(userID, dateIndex, ctxLogger).
				Times(1).
				Return(errors.New("database error"))

			err := service.DeleteGymAttendance(userID, date, ctxLogger)
			Expect(err).ToNot(BeNil())
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
