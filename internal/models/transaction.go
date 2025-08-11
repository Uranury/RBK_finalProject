package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Withdraw TransactionType = "withdraw"
	Deposit  TransactionType = "deposit"
)

type Transaction struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	Amount        float64         `json:"amount" db:"amount"`
	Type          TransactionType `json:"type" db:"type"`
	BalanceBefore float64         `json:"balance_before" db:"balance_before"`
	BalanceAfter  float64         `json:"balance_after" db:"balance_after"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
