package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/H0wZy/authsys/internal/dto"
	"github.com/H0wZy/authsys/internal/jwt"
	"github.com/H0wZy/authsys/internal/model"
	"github.com/H0wZy/authsys/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Login(ctx context.Context, dto *dto.Auth) (string, error)
	Logout(ctx context.Context, user *model.User) error
}

type authService struct {
	repo repository.UserRepository
	jwtm jwt.JwtManager
}

func (s *authService) Login(ctx context.Context, dto *dto.Auth) (string, error) {
	var user *model.User
	var err error

	if strings.Contains(dto.Login, "@") {
		user, err = s.repo.FindByEmail(ctx, dto.Login)
	} else {
		user, err = s.repo.FindByUsername(ctx, dto.Login)
	}

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("error finding user: %w", err)
	}

	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); bcryptErr != nil {
		return "", ErrInvalidCredentials
	}

	return s.jwtm.Generate(user.ID, user.Email, user.Username)
}

func (s *authService) Logout(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}

func NewAuthService(repo repository.UserRepository, jwtm jwt.JwtManager) AuthService {
	return &authService{repo: repo, jwtm: jwtm}
}
