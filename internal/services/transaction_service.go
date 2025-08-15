package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/repositories/transaction"
	"github.com/Uranury/RBK_finalProject/internal/repositories/user"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TransactionService struct {
	transactionRepo transaction.Repository
	userRepo        user.Repository
	db              *sqlx.DB
	logger          *slog.Logger
}

func NewTransactionService(transactionRepo transaction.Repository, userRepo user.Repository, db *sqlx.DB, logger *slog.Logger) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		db:              db,
		logger:          logger,
	}
}

// Withdraw handles withdrawing money from user's balance
func (s *TransactionService) Withdraw(ctx context.Context, userID uuid.UUID, amount float64) (*models.Transaction, error) {
	s.logger.Info("starting withdrawal", "user_id", userID, "amount", amount)

	// Validate amount
	if amount <= 0 {
		return nil, apperrors.NewValidationError("Withdrawal amount must be greater than zero")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to start transaction", "error", err)
		return nil, apperrors.NewInternalError("Failed to process withdrawal", err)
	}
	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil {
			s.logger.Error("failed to rollback transaction", "error", err)
		}
	}(tx)

	// Get current user balance
	usr, err := s.userRepo.GetUserByIdForUpdate(ctx, tx, userID)
	if err != nil {
		s.logger.Error("failed to get user for update", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process withdrawal", err)
	}
	if usr == nil {
		return nil, apperrors.NewNotFoundError("User not found")
	}

	balanceBefore := usr.Balance
	if balanceBefore < amount {
		s.logger.Warn("insufficient funds for withdrawal", "user_id", userID, "balance", balanceBefore, "requested", amount)
		return nil, apperrors.NewValidationError("Insufficient funds for withdrawal")
	}

	// Update user balance
	newBalance := balanceBefore - amount
	if err := s.userRepo.UpdateBalance(ctx, tx, userID, newBalance); err != nil {
		s.logger.Error("failed to update user balance", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process withdrawal", err)
	}

	// Create transaction record
	now := time.Now()
	trnsc := &models.Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        amount,
		Type:          models.Withdraw,
		BalanceBefore: balanceBefore,
		BalanceAfter:  newBalance,
		CreatedAt:     now,
	}

	if err := s.transactionRepo.Create(ctx, trnsc); err != nil {
		s.logger.Error("failed to create transaction record", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process withdrawal", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process withdrawal", err)
	}

	s.logger.Info("withdrawal completed successfully",
		"user_id", userID,
		"amount", amount,
		"balance_before", balanceBefore,
		"balance_after", newBalance)

	return trnsc, nil
}

// Deposit handles depositing money to user's balance
func (s *TransactionService) Deposit(ctx context.Context, userID uuid.UUID, amount float64) (*models.Transaction, error) {
	s.logger.Info("starting deposit", "user_id", userID, "amount", amount)

	// Validate amount
	if amount <= 0 {
		return nil, apperrors.NewValidationError("Deposit amount must be greater than zero")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to start transaction", "error", err)
		return nil, apperrors.NewInternalError("Failed to process deposit", err)
	}
	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil {
			s.logger.Error("failed to rollback transaction", "error", err)
		}
	}(tx)

	// Get current user balance
	usr, err := s.userRepo.GetUserByIdForUpdate(ctx, tx, userID)
	if err != nil {
		s.logger.Error("failed to get user for update", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process deposit", err)
	}
	if usr == nil {
		return nil, apperrors.NewNotFoundError("User not found")
	}

	balanceBefore := usr.Balance
	newBalance := balanceBefore + amount

	// Update user balance
	if err := s.userRepo.UpdateBalance(ctx, tx, userID, newBalance); err != nil {
		s.logger.Error("failed to update user balance", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process deposit", err)
	}

	// Create transaction record
	now := time.Now()
	trnsc := &models.Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        amount,
		Type:          models.Deposit,
		BalanceBefore: balanceBefore,
		BalanceAfter:  newBalance,
		CreatedAt:     now,
	}

	if err := s.transactionRepo.Create(ctx, trnsc); err != nil {
		s.logger.Error("failed to create transaction record", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process deposit", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", "error", err, "user_id", userID)
		return nil, apperrors.NewInternalError("Failed to process deposit", err)
	}

	s.logger.Info("deposit completed successfully",
		"user_id", userID,
		"amount", amount,
		"balance_before", balanceBefore,
		"balance_after", newBalance)

	return trnsc, nil
}

// GetUserTransactions returns transaction history for a user
func (s *TransactionService) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]*models.Transaction, error) {
	transactions, err := s.transactionRepo.GetUserTransactions(ctx, userID)
	if err != nil {
		return nil, apperrors.NewInternalError("Failed to retrieve transaction history", err)
	}
	return transactions, nil
}
