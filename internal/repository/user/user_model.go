package user

import (
	badgeModelDB "gym-badges-api/internal/repository/badge"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	ID          string       `gorm:"primary_key;not null" json:"user_id"`
	BodyFat     float32      `gorm:"not null;type:decimal(5,2)" json:"body_fat"`
	CurrentWeek pq.BoolArray `gorm:"not null;type:bool[]" json:"current_week"`
	Email       string       `gorm:"not null;unique" json:"email"`
	Experience  int64        `gorm:"not null" json:"experience"`
	Image       string       `gorm:"null" json:"image"`
	Name        string       `gorm:"not null" json:"name"`
	Password    string       `gorm:"not null" json:"password"`
	Streak      int32        `gorm:"not null" json:"streak"`
	WeeklyGoal  int32        `gorm:"not null" json:"weekly_goal"`
	Weight      float32      `gorm:"not null;type:decimal(5,2)" json:"weight"`

	GymAttendance []GymAttendance       `gorm:"constraint:OnDelete:CASCADE"`
	FatHistory    []FatHistory          `gorm:"constraint:OnDelete:CASCADE"`
	WeightHistory []WeightHistory       `gorm:"constraint:OnDelete:CASCADE"`
	Friends       []*User               `gorm:"many2many:user_friends;constraint:OnDelete:CASCADE"`
	Badges        []*badgeModelDB.Badge `gorm:"many2many:user_badges;constraint:OnDelete:CASCADE"`
	TopFeats      []*badgeModelDB.Badge `gorm:"many2many:user_top_feats;constraint:OnDelete:CASCADE"`
	Preferences   []Preference          `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"null" json:"created_at"`
	UpdatedAt time.Time `gorm:"null" json:"updated_at"`
	DeletedAt time.Time `gorm:"null" json:"deleted_at"`
}

type GymAttendance struct {
	UserID string    `gorm:"primary_key;not null"`
	Date   time.Time `gorm:"primary_key;not null"`

	CreatedAt time.Time `gorm:"null" json:"created_at"`
	UpdatedAt time.Time `gorm:"null" json:"updated_at"`
	DeletedAt time.Time `gorm:"null" json:"deleted_at"`
}

type FatHistory struct {
	UserID string    `gorm:"primary_key;not null"`
	Date   time.Time `gorm:"primary_key;not null"`
	Fat    float32   `gorm:"not null;type:decimal(5,2)"`

	CreatedAt time.Time `gorm:"null" json:"created_at"`
	UpdatedAt time.Time `gorm:"null" json:"updated_at"`
	DeletedAt time.Time `gorm:"null" json:"deleted_at"`
}

type WeightHistory struct {
	UserID string    `gorm:"primary_key;not null"`
	Date   time.Time `gorm:"primary_key;not null"`
	Weight float32   `gorm:"not null;type:decimal(5,2)"`

	CreatedAt time.Time `gorm:"null" json:"created_at"`
	UpdatedAt time.Time `gorm:"null" json:"updated_at"`
	DeletedAt time.Time `gorm:"null" json:"deleted_at"`
}

type Preference struct {
	gorm.Model
	PreferenceInfo PreferenceInfo `gorm:"not null"`
	On             bool           `gorm:"not null"`

	UserID string `gorm:"not null"`
}

type PreferenceInfo struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`

	PreferenceID uint `gorm:"not null"`
}
