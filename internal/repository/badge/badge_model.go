package badge_dao

import (
	"gorm.io/gorm"
)

type Badge struct {
	gorm.Model
	Name          string `gorm:"not null"`
	Description   string `gorm:"not null"`
	Image         string `gorm:"not null"`
	ParentBadgeID uint   `gorm:"null"`
	ParentBadge   *Badge `gorm:"null"`
}
