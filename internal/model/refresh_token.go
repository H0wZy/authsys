package model

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"not null"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash []byte    `gorm:"type:bytea;uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
}
