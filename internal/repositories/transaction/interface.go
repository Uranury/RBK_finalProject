package transaction

import (
	"context"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error)
}
