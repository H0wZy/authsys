package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Name         string     `gorm:"type:varchar(200);not null" json:"full_name"`
	Password     string     `gorm:"not null" json:"-"`
	Phone        string     `gorm:"type:varchar(20);not null" json:"phone"`
	BirthDate    *time.Time `json:"birth_date,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	LastLogoutAt *time.Time `json:"last_logout_at,omitempty"`
	IsUserActive bool       `gorm:"default:true" json:"is_user_active"`
}
