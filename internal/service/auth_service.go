package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/H0wZy/authsys/internal/dto"
	"github.com/H0wZy/authsys/internal/jwt"
	"github.com/H0wZy/authsys/internal/model"
	"github.com/H0wZy/authsys/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenTTL   = 7 * 24 * time.Hour
	refreshTokenBytes = 32
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthService interface {
	Login(ctx context.Context, input *dto.Auth) (*dto.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	repo   repository.UserRepository
	tokens repository.RefreshTokenRepository
	jwtm   jwt.Manager
}

func (s *authService) Login(ctx context.Context, input *dto.Auth) (*dto.AuthResponse, error) {
	var user *model.User
	var err error

	if strings.Contains(input.Login, "@") {
		user, err = s.repo.FindByEmail(ctx, input.Login)
	} else {
		user, err = s.repo.FindByUsername(ctx, input.Login)
	}

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); bcryptErr != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issuePair(ctx, user)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	stored, err := s.tokens.FindByHash(ctx, hashToken(refreshToken))
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("error finding refresh token: %w", err)
	}

	now := time.Now()

	// Um token já revogado voltando significa que alguém guardou uma cópia.
	// Não dá para saber se quem está apresentando é o dono ou o atacante,
	// então derruba todas as sessões e obriga login de novo.
	if stored.RevokedAt != nil {
		if err := s.tokens.RevokeAllForUser(ctx, stored.UserID, now); err != nil {
			return nil, fmt.Errorf("error revoking token family: %w", err)
		}
		return nil, ErrInvalidRefreshToken
	}

	if !now.Before(stored.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.repo.FindByID(ctx, stored.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// Rotação: queima o token apresentado antes de emitir o próximo.
	if err := s.tokens.Revoke(ctx, stored.ID, now); err != nil {
		return nil, fmt.Errorf("error rotating refresh token: %w", err)
	}

	return s.issuePair(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	stored, err := s.tokens.FindByHash(ctx, hashToken(refreshToken))
	if err != nil {
		// Token inexistente já está, para todos os efeitos, deslogado.
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil
		}
		return fmt.Errorf("error finding refresh token: %w", err)
	}

	if err := s.tokens.Revoke(ctx, stored.ID, time.Now()); err != nil {
		return fmt.Errorf("error revoking refresh token: %w", err)
	}

	return nil
}

// issuePair emite um access token novo e um refresh token novo, gravando
// apenas o hash do segundo. É o único lugar que constrói um dto.AuthResponse.
func (s *authService) issuePair(ctx context.Context, user *model.User) (*dto.AuthResponse, error) {
	access, err := s.jwtm.Generate(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	raw, hash, err := newRefreshToken()
	if err != nil {
		return nil, err
	}

	record := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}

	if err := s.tokens.Create(ctx, record); err != nil {
		return nil, fmt.Errorf("error storing refresh token: %w", err)
	}

	return &dto.AuthResponse{AccessToken: access, RefreshToken: raw}, nil
}

// newRefreshToken devolve o token em claro (que só o cliente verá) e o hash
// que vai para o banco.
func newRefreshToken() (string, []byte, error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	raw := base64.RawURLEncoding.EncodeToString(b)

	return raw, hashToken(raw), nil
}

// SHA-256 basta porque o token tem 256 bits de entropia: não há dicionário
// para atacar. bcrypt aqui só custaria tempo e impediria a busca por índice.
func hashToken(raw string) []byte {
	sum := sha256.Sum256([]byte(raw))
	return sum[:]
}

func NewAuthService(
	repo repository.UserRepository,
	tokens repository.RefreshTokenRepository,
	jwtm jwt.Manager,
) AuthService {
	return &authService{repo: repo, tokens: tokens, jwtm: jwtm}
}
