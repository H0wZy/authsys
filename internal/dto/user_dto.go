package dto

import (
	"time"

	"github.com/H0wZy/authsys/internal/model"
)

type CreateUser struct {
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Password  string    `json:"password" binding:"required,min=8"`
	Phone     string    `json:"phone" binding:"required"`
	BirthDate time.Time `json:"birth_date"`
}

type UserResponse struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(u *model.User) UserResponse {
	return UserResponse{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Phone:     u.Phone,
		BirthDate: u.BirthDate,
		CreatedAt: u.CreatedAt,
	}
}

type UpdateUser struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type UpdatePassword struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}
