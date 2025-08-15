package skin

import (
	"context"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, skin *models.Skin) error
	GetSkin(ctx context.Context, id uuid.UUID) (*models.Skin, error)
	GetUserSkins(ctx context.Context, userID uuid.UUID) ([]*models.Skin, error)
	GetAvailableSkins(ctx context.Context) ([]*models.Skin, error)
	GetSkinsForUpdate(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID) ([]*models.Skin, error)
	GetSkinsForSellUpdate(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID) ([]*models.Skin, error)
	UpdateOwnership(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID, newOwnerID uuid.UUID) error
	UpdatePrice(ctx context.Context, tx *sqlx.Tx, skinID uuid.UUID, price float64) error
	UpdateForSale(ctx context.Context, tx *sqlx.Tx, skinID uuid.UUID, price float64, available bool) error
	UpdateAvailability(ctx context.Context, tx *sqlx.Tx, skinID uuid.UUID, available bool) error
}
