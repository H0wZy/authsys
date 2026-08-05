package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/H0wZy/authsys/internal/dto"
	"github.com/H0wZy/authsys/internal/errorlist"
	"github.com/H0wZy/authsys/internal/model"
	"github.com/H0wZy/authsys/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(ctx context.Context, input *dto.CreateUser) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func (s *userService) Create(ctx context.Context, input *dto.CreateUser) (*model.User, error) {
	if input.BirthDate.After(time.Now()) {
		return nil, errorlist.InvalidBirthDate
	}

	if strings.Contains(input.Username, "@") {
		return nil, errorlist.UsernameCantContainAt
	}

	_, err := s.repo.FindByEmail(ctx, input.Email)
	if err == nil {
		return nil, errorlist.EmailAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking email: %w", err)
	}

	_, err = s.repo.FindByUsername(ctx, input.Username)
	if err == nil {
		return nil, errorlist.UsernameAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error while hashing password: %w", err)
	}

	user := &model.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hash),
		Phone:     input.Phone,
		BirthDate: input.BirthDate,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
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
