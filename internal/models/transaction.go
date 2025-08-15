package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Withdraw TransactionType = "withdraw"
	Deposit  TransactionType = "deposit"
	Purchase TransactionType = "purchase"
	Sale     TransactionType = "sale"
)

type Transaction struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	Amount        float64         `json:"amount" db:"amount"`
	Type          TransactionType `json:"type" db:"type"`
	BalanceBefore float64         `json:"balance_before" db:"balance_before"`
	BalanceAfter  float64         `json:"balance_after" db:"balance_after"`

	SkinID         *uuid.UUID `json:"skin_id,omitempty" db:"skin_id"`
	OrderID        *uuid.UUID `json:"order_id,omitempty" db:"order_id"`
	CounterpartyID *uuid.UUID `json:"counterparty_id,omitempty" db:"counterparty_id"`
	Description    *string    `json:"description,omitempty" db:"description"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
