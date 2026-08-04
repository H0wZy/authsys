package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/H0wZy/authsys/internal/dto"
	"github.com/H0wZy/authsys/internal/model"
	"github.com/H0wZy/authsys/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(ctx context.Context, dto *dto.Auth) (string, error)
	Logout(ctx context.Context, user *model.User) error
}

type UserService interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type service struct {
	repo repository.UserRepository
}

func (s *service) Create(ctx context.Context, user *model.User) error {

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

func (s *service) Delete(ctx context.Context, id uint) error {
	panic("unimplemented")
}

func (s *service) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("unimplemented")
}

func (s *service) FindByID(ctx context.Context, id uint) (*model.User, error) {
	panic("unimplemented")
}

func (s *service) List(ctx context.Context) ([]*model.User, error) {
	panic("unimplemented")
}

func (s *service) Update(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}

func (s *service) Login(ctx context.Context, dto *dto.Auth) (string, error) {
	user, err := s.repo.FindByEmail(ctx, dto.Login)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	tokenString, err := s.generateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		return "", fmt.Errorf("error while generating token: %w", err)
	}

	return tokenString, nil
}

func (s *service) Logout(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}

func (s *service) generateJWT(id uint, email string, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      id,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &service{repo: repo}
}

func NewUserService(repo repository.UserRepository) UserService {
	return &service{repo: repo}
}
