package user

type User struct {
	UserID      string  `gorm:"primary_key;not null" json:"user_id"`
	BodyFat     float32 `gorm:"null" json:"body_fat"`
	CurrentWeek []bool  `gorm:"null;type:bool[]" json:"current_week"`
	Email       string  `gorm:"not null;unique" json:"email"`
	Experience  int64   `gorm:"null" json:"experience"`
	Image       string  `gorm:"null" json:"image"`
	Name        string  `gorm:"null" json:"name"`
	Password    string  `gorm:"null" json:"password"`
	Streak      int32   `gorm:"null" json:"streak"`
	Weight      float32 `gorm:"null" json:"weight"`
	CreatedAt   string  `gorm:"null" json:"created_at"`
	UpdatedAt   string  `gorm:"null" json:"updated_at"`
}
