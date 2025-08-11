package transaction

import (
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	WriteTransaction(ctx context.Context, tx *sqlx.Tx, transaction *models.Transaction) error
}
