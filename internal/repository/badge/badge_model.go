package badge_dao

type Badge struct {
	ID            int16  `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Description   string `gorm:"not null"`
	Image         string `gorm:"not null"`
	ParentBadgeID int16  `gorm:"null"`
	ParentBadge   *Badge `gorm:"null"`
}
