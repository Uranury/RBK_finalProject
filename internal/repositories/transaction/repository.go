package transaction

import (
	"context"
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

func (r *repository) WriteTransaction(ctx context.Context, tx *sqlx.Tx, transaction *models.Transaction) error {
	query := `
    INSERT INTO transactions(id, sender_id, receiver_id, amount, type, created_at) 
    VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.ExecContext(ctx, query, transaction.ID, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.Type, transaction.CreatedAt)
	return err
}

func (r *repository) IncreaseBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, amount int64) error {
	query := `
	UPDATE users SET balance = balance + $1 WHERE id = $2
	`
	_, err := tx.ExecContext(ctx, query, amount, userID)
	return err
}
