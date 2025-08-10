package user

import (
	"context"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetBalance(ctx context.Context, userID uuid.UUID) (float64, error)
	UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, newBalance float64) error
	Create(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByIdForUpdate(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) (*models.User, error)
}
