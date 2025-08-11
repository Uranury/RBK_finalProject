package transaction

import (
	"context"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, tx *sqlx.Tx, transaction *models.Transaction) error
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error)
}
