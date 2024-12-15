package user

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type IUserDAO interface {
	GetUser(userID string, ctxLog *log.Entry) (*User, error)
	GetUserByEmail(email string, ctxLog *log.Entry) (*User, error)
	CreateUser(user *User, ctxLog *log.Entry) error
	EditUserInfo(userID string, newUserInfo *User, ctxLog *log.Entry) (*User, error)
	AddDayToCurrentWeek(userID string, dayIndex int, ctxLog *log.Entry) error
	DeleteDayFromCurrentWeek(userID string, dayIndex int, ctxLog *log.Entry) error

	GetUserWithWeightHistory(userID string, months int32, ctxLog *log.Entry) (*User, error)
	AddWeight(userID string, weight float32, date time.Time, ctxLog *log.Entry) error

	GetUserWithFatHistory(userID string, months int32, ctxLog *log.Entry) (*User, error)
	AddBodyFat(userID string, bodyFat float32, date time.Time, ctxLog *log.Entry) error

	GetUserWithAttendance(userID string, year int32, month int32, ctxLog *log.Entry) (*User, error)
	AddGymAttendance(userID string, date time.Time, ctxLog *log.Entry) error
	DeleteGymAttendance(userID string, date time.Time, ctxLog *log.Entry) error

	GetUserWithFriends(userID string, offset int32, size int32, ctxLog *log.Entry) (*User, error)
	GetUserWithBadges(userID string, ctxLog *log.Entry) (*User, error)
	AddFriend(userID string, friendID string, ctxLog *log.Entry) (*User, error)
	DeleteFriend(userID string, friendID string, ctxLog *log.Entry) error

	AddExperience(userID string, exp int64, ctxLog *log.Entry) error
}
