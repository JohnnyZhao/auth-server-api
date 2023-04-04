package domain

import (
	"time"
)

type User struct {
	ID             uint      `json:"-"`
	UserID         string    `gorm:"type:varchar(24);uniqueIndex;not null" json:"user_id"`
	HashedPassword string    `gorm:"type:varchar(256);not null" json:"-"`
	Nickname       string    `gorm:"type:varchar(32);" json:"nickname"`
	Comment        *string   `gorm:"type:varchar(128);" json:"comment,omitempty"`
	CreatedAt      time.Time `gorm:"not null" json:"-"`
	UpdatedAt      time.Time `gorm:"not null" json:"-"`
}
