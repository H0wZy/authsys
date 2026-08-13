package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/H0wZy/authsys/internal/model"
	"gorm.io/gorm"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	FindByHash(ctx context.Context, hash []byte) (*model.RefreshToken, error)
	Revoke(ctx context.Context, id uint, at time.Time) error
	RevokeAllForUser(ctx context.Context, userID uint, at time.Time) error
}

type refreshTokenRepository struct{ db *gorm.DB }

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("error creating refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) FindByHash(ctx context.Context, hash []byte) (*model.RefreshToken, error) {
	var token model.RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("error finding refresh token: %w", err)
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uint, at time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", at).Error; err != nil {
		return fmt.Errorf("error revoking refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uint, at time.Time) error {
	if err := r.db.WithContext(ctx).
		Model(&model.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", at).Error; err != nil {
		return fmt.Errorf("error revoking tokens for user %d: %w", userID, err)
	}
	return nil
}
