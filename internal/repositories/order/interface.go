package order

import (
	"context"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, tx *sqlx.Tx, order *models.Order) error
	CreateOrderItem(ctx context.Context, tx *sqlx.Tx, orderItem *models.OrderItem) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetOrderByIDForUpdate(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*models.Order, error)
	UpdateStatus(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, status models.OrderStatus) error
}
