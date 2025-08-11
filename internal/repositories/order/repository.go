package order

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

func (r *repository) Create(ctx context.Context, tx *sqlx.Tx, order *models.Order) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO orders (id, user_id, total_amount, status, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6)`,
		order.ID, order.UserID, order.TotalAmount, order.Status, order.CreatedAt, order.UpdatedAt)

	return err
}

func (r *repository) CreateOrderItem(ctx context.Context, tx *sqlx.Tx, orderItem *models.OrderItem) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO order_items (id, order_id, skin_id, price, created_at)
				VALUES ($1, $2, $3, $4, $5)`,
		orderItem.ID, orderItem.OrderID, orderItem.SkinID, orderItem.Price, orderItem.CreatedAt)
	return err
}

func (r *repository) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	err := r.db.GetContext(ctx, order, "SELECT * FROM orders WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return order, nil
}

func (r *repository) GetOrderByIDForUpdate(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*models.Order, error) {
	order := &models.Order{}
	query := `
        SELECT * 
        FROM orders 
        WHERE id = $1 
        FOR UPDATE
    `
	err := tx.GetContext(ctx, order, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}

func (r *repository) UpdateStatus(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, status models.OrderStatus) error {
	_, err := tx.ExecContext(ctx, `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}
