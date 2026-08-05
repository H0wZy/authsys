package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/H0wZy/authsys/internal/model"
	"github.com/H0wZy/authsys/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) Create(ctx context.Context, user *model.User) error {

	if strings.Contains(user.Username, "@") {
		return errors.New("username cannot contain '@'")
	}

	_, err := s.repo.FindByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("email already in use")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking email: %w", err)
	}

	_, err = s.repo.FindByUsername(ctx, user.Username)
	if err == nil {
		return errors.New("username already in use")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error while hashing password: %w", err)
	}

	user.Password = string(hash)

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("error while creating user: %w", err)
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	panic("unimplemented")
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("unimplemented")
}

func (s *userService) FindByID(ctx context.Context, id uint) (*model.User, error) {
	panic("unimplemented")
}

func (s *userService) List(ctx context.Context) ([]*model.User, error) {
	panic("unimplemented")
}

func (s *userService) Update(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
