package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	UserID      uuid.UUID   `json:"userId" db:"user_id"`
	TotalAmount float64     `json:"totalAmount" db:"total_amount"`
	Status      OrderStatus `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OrderID   uuid.UUID `json:"orderId" db:"order_id"`
	SkinID    uuid.UUID `json:"skinId" db:"skin_id"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
