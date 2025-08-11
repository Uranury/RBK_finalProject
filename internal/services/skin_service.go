package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/repositories/skin"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
)

type Skin struct {
	repo   skin.Repository
	logger *slog.Logger
}

func NewSkin(repo skin.Repository, logger *slog.Logger) *Skin {
	return &Skin{repo: repo, logger: logger}
}

// GetAllGuns returns all available guns in the system
func (s *Skin) GetAllGuns() []models.Gun {
	return []models.Gun{
		// Pistols
		models.AK47, models.M4A4, models.M4A1S, models.DesertEagle, models.USPS,
		models.Glock18, models.P250, models.Tec9, models.CZ75,
		// Rifles
		models.AWP, models.SSG08, models.SCAR20, models.G3SG1,
		// SMGs
		models.MP9, models.MAC10, models.MP7, models.P90, models.UMP45, models.PPBizon,
		// Shotguns
		models.Nova, models.XM1014, models.MAG7, models.SawedOff,
		// Machine Guns
		models.M249, models.Negev,
		// Knives
		models.Karambit, models.Butterfly, models.M9Bayonet, models.Bayonet,
		models.FlipKnife, models.GutKnife, models.Huntsman, models.ShadowDaggers,
		// Other
		models.Falchion, models.Bowie, models.Navaja, models.Stiletto,
		models.Ursus, models.Nomad, models.Paracord, models.Survival, models.Classic,
	}
}

// GetAllWears returns all available wear levels in the system
func (s *Skin) GetAllWears() []models.Wear {
	return []models.Wear{
		models.FactoryNew,
		models.MinimalWear,
		models.FieldTested,
		models.WellWorn,
		models.BattleScarred,
	}
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
	if skin.Condition < 0 || skin.Condition > 1 {
		s.logger.Warn("skin creation failed: invalid condition", "condition", skin.Condition)
		return nil, apperrors.NewValidationError("Skin condition must be between 0 and 1")
	}

	// Set default gun if not provided
	if skin.Gun == "" {
		skin.Gun = models.AK47
	}

	// Automatically calculate wear based on condition
	skin.Wear = models.GetWearFromCondition(skin.Condition)

	now := time.Now()
	skin.ID = uuid.New()
	skin.CreatedAt = now
	skin.UpdatedAt = now
	skin.Available = true // New skins are available by default
	skin.OwnerID = nil

	if err := s.repo.Create(ctx, skin); err != nil {
		s.logger.Error("failed to create skin in repository", "error", err)
		return nil, apperrors.NewInternalError("Failed to create skin", err)
	}

	s.logger.Info("skin created successfully", "skin_id", skin.ID)
	return skin, nil
}

func (s *Skin) GetSkinByID(ctx context.Context, id uuid.UUID) (*models.Skin, error) {
	sk, err := s.repo.GetSkin(ctx, id)
	if err != nil {
		s.logger.Error("failed to get skin in repository", "error", err)
		return nil, apperrors.NewInternalError("Failed to retrieve skin", err)
	}
	if sk == nil {
		return nil, apperrors.NewNotFoundError("Skin not found")
	}
	return sk, nil
}

func (s *Skin) GetAvailableSkins(ctx context.Context) ([]*models.Skin, error) {
	sks, err := s.repo.GetAvailableSkins(ctx)
	if err != nil {
		s.logger.Error("failed to get available skins in repository", "error", err)
		return nil, apperrors.NewInternalError("Failed to retrieve available skins", err)
	}
	if len(sks) == 0 {
		return []*models.Skin{}, nil
	}
	return sks, nil
}

func (s *Skin) GetUserSkins(ctx context.Context, userID uuid.UUID) ([]*models.Skin, error) {
	sks, err := s.repo.GetUserSkins(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get user skins in repository", "error", err)
		return nil, apperrors.NewInternalError("Failed to retrieve user skins", err)
	}
	if len(sks) == 0 {
		return []*models.Skin{}, nil
	}
	return sks, nil
}
