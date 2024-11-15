package user

import "github.com/lib/pq"

type User struct {
	UserID      string       `gorm:"primary_key;not null" json:"user_id"`
	BodyFat     float32      `gorm:"not null" json:"body_fat"`
	CurrentWeek pq.BoolArray `gorm:"not null;type:bool[]" json:"current_week"`
	Email       string       `gorm:"not null;unique" json:"email"`
	Experience  int64        `gorm:"not null" json:"experience"`
	Image       string       `gorm:"null" json:"image"`
	Name        string       `gorm:"not null" json:"name"`
	Password    string       `gorm:"not null" json:"password"`
	Streak      int32        `gorm:"not null" json:"streak"`
	Weight      float32      `gorm:"not null" json:"weight"`
	CreatedAt   string       `gorm:"null" json:"created_at"`
	UpdatedAt   string       `gorm:"null" json:"updated_at"`
}
