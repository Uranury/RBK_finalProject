package skin

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateSkin(ctx context.Context, skin *models.Skin) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO skins (id, owner_id, name, rarity, float, price, image, available, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		skin.ID, skin.OwnerID, skin.Name, skin.Rarity, skin.Condition, skin.Price, skin.Image, skin.Available, skin.CreatedAt, skin.UpdatedAt,
	)
	return err // Simplified return
}

func (r *repository) GetSkin(ctx context.Context, id uuid.UUID) (*models.Skin, error) {
	var skin models.Skin
	err := r.db.GetContext(ctx, &skin, "SELECT * FROM skins WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &skin, nil
}

func (r *repository) GetUserSkins(ctx context.Context, userID uuid.UUID) ([]*models.Skin, error) {
	var skins []*models.Skin
	err := r.db.SelectContext(ctx, &skins, "SELECT * FROM skins WHERE owner_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Skin{}, nil // Return empty slice, not nil
		}
		return nil, err
	}
	return skins, nil
}

func (r *repository) GetAvailableSkins(ctx context.Context) ([]*models.Skin, error) {
	var skins []*models.Skin
	err := r.db.SelectContext(ctx, &skins,
		"SELECT * FROM skins WHERE available = true AND owner_id IS NULL ORDER BY created_at DESC")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Skin{}, nil // Return empty slice, not nil
		}
		return nil, err
	}
	return skins, nil
}

// GetSkinsForUpdate You'll need these methods for the marketplace transaction:
func (r *repository) GetSkinsForUpdate(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID) ([]*models.Skin, error) {
	query, args, err := sqlx.In(
		"SELECT * FROM skins WHERE id IN (?) AND available = true AND owner_id IS NULL FOR UPDATE",
		skinIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var skins []*models.Skin
	err = tx.SelectContext(ctx, &skins, query, args...)
	return skins, err
}

func (r *repository) UpdateOwnership(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID, newOwnerID uuid.UUID) error {
	query, args, err := sqlx.In(
		"UPDATE skins SET owner_id = ?, available = false, updated_at = NOW() WHERE id IN (?)",
		newOwnerID, skinIDs)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)

	_, err = tx.ExecContext(ctx, query, args...)
	return err
}
