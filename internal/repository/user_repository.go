package repository

import (
	"context"
	"fmt"

	"github.com/H0wZy/authsys/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := u.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (u *userRepository) Delete(ctx context.Context, id uint) error {
	result := u.db.WithContext(ctx).Delete(&model.User{}, id)

	if result.Error != nil {
		return fmt.Errorf("error deleting user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("error deleting user: %w", gorm.ErrRecordNotFound)
	}

	return nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return &user, nil
}

func (u *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	if err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("error finding user by username: %w", err)
	}

	return &user, nil
}

func (u *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("error finding user by id: %w", err)
	}

	return &user, nil
}

func (u *userRepository) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User

	if err := u.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (u *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := u.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
