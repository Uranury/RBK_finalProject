package transaction

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

func (r *repository) Create(ctx context.Context, transaction *models.Transaction) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO transaction_history (id, user_id, amount, type, balance_before, balance_after, skin_id, 
                                 order_id, counterparty_id, description, created_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		transaction.ID, transaction.UserID, transaction.Amount, transaction.Type,
		transaction.BalanceBefore, transaction.BalanceAfter, transaction.SkinID,
		transaction.OrderID, transaction.CounterpartyID, transaction.Description,
		transaction.CreatedAt)
	return err
}

func (r *repository) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := r.db.SelectContext(ctx, &transactions,
		"SELECT * FROM transaction_history WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Transaction{}, nil
		}
		return nil, err
	}
	return transactions, nil
}
