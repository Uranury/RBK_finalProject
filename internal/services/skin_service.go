package services

import (
	"context"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/repositories/skin"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Skin struct {
	repo   skin.Repository
	logger *slog.Logger
}

func NewSkin(repo skin.Repository, logger *slog.Logger) *Skin {
	return &Skin{repo: repo, logger: logger}
}

func (s *Skin) CreateSkin(ctx context.Context, skin *models.Skin) (*models.Skin, error) {
	if skin.Name == "" {
		s.logger.Warn("skin creation failed: missing name", "skin", skin)
		return nil, apperrors.NewValidationError("skin name is required")
	}
	if skin.Rarity == "" {
		s.logger.Warn("skin creation failed: missing rarity", "skin", skin)
		return nil, apperrors.NewValidationError("skin rarity is required")
	}
	if skin.Price <= 0 {
		s.logger.Warn("skin creation failed: invalid price", "price", skin.Price)
		return nil, apperrors.NewValidationError("skin price can't be negative or zero")
	}
	if skin.Condition > 1 || skin.Condition < 0 {
		s.logger.Warn("skin creation failed: invalid condition", "condition", skin.Condition)
		return nil, apperrors.NewValidationError("skin condition should be between 0 and 1")
	}

	now := time.Now()
	skin.ID = uuid.New()
	skin.CreatedAt = now
	skin.UpdatedAt = now
	skin.Available = true // New skins are available by default
	skin.OwnerID = nil

	if err := s.repo.Create(ctx, skin); err != nil {
		s.logger.Error("failed to create skin in repository", "error", err)
		return nil, apperrors.WrapInternal(err, "failed to create skin")
	}

	s.logger.Info("skin created successfully", "skin_id", skin.ID)
	return skin, nil
}

func (s *Skin) GetSkinByID(ctx context.Context, id uuid.UUID) (*models.Skin, error) {
	sk, err := s.repo.GetSkin(ctx, id)
	if err != nil {
		s.logger.Error("failed to get skin in repository", "error", err)
		return nil, apperrors.WrapInternal(err, "failed to get skin")
	}
	if sk == nil {
		return nil, apperrors.NewNotFoundError("skin does not exist")
	}
	return sk, nil
}

func (s *Skin) GetAvailableSkins(ctx context.Context) ([]*models.Skin, error) {
	sks, err := s.repo.GetAvailableSkins(ctx)
	if err != nil {
		s.logger.Error("failed to get skin in repository", "error", err)
		return nil, apperrors.WrapInternal(err, "failed to get skins")
	}
	if len(sks) == 0 {
		return []*models.Skin{}, nil
	}
	return sks, nil
}

func (s *Skin) GetUserSkins(ctx context.Context, userID uuid.UUID) ([]*models.Skin, error) {
	sks, err := s.repo.GetUserSkins(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get skin in repository", "error", err)
		return nil, apperrors.WrapInternal(err, "failed to get skins")
	}
	if len(sks) == 0 {
		return []*models.Skin{}, nil
	}
	return sks, nil
}
