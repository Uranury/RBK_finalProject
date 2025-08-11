package models

import (
	"github.com/google/uuid"
	"time"
)

type TransactionType string

const (
	Withdraw TransactionType = "withdraw"
	Deposit  TransactionType = "deposit"
)

type Transaction struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	SenderID   uuid.UUID       `json:"sender_id" db:"sender_id"`
	ReceiverID uuid.UUID       `json:"receiver_id" db:"receiver_id"`
	Amount     int64           `json:"amount" db:"amount"`
	Type       TransactionType `json:"type" db:"type"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}
