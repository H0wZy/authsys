package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	Password  string     `gorm:"type:varchar(60);not null" json:"-"`
	Phone     string     `gorm:"type:varchar(20);not null" json:"phone"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	Account
}

type Account struct {
	IsOnline            bool       `json:"is_online"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	LastLogoutAt        *time.Time `json:"last_logout_at,omitempty"`
	FailedLoginAttempts int        `gorm:"default:0" json:"-"`
	LockedUntil         *time.Time `json:"-"`
	IsAccountDisabled   bool       `gorm:"default:true" json:"is_account_disabled"`
}
