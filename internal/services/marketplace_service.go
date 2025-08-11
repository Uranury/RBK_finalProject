package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/repositories/order"
	"github.com/Uranury/RBK_finalProject/internal/repositories/skin"
	"github.com/Uranury/RBK_finalProject/internal/repositories/user"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
)

// TODO: Write orderRepo methods
// TODO: Get rid of queue stuff in purchase, get the MVP first

type MarketplaceService struct {
	skinRepo   skin.Repository
	orderRepo  order.Repository
	userRepo   user.Repository
	emailQueue *asynq.Client
	db         *sqlx.DB
	logger     *slog.Logger
}

func NewMarketplaceService(skinRepo skin.Repository,
	orderRepo order.Repository,
	userRepo user.Repository,
	emailQueue *asynq.Client,
	db *sqlx.DB,
	logger *slog.Logger) *MarketplaceService {
	return &MarketplaceService{skinRepo, orderRepo, userRepo, emailQueue, db, logger}
}

// PurchaseSkin handles buying a single skin
func (s *MarketplaceService) PurchaseSkin(ctx context.Context, userID uuid.UUID, skinID uuid.UUID) (*models.Order, error) {
	s.logger.Info("starting skin purchase", "user_id", userID, "skin_id", skinID)

	// Start database transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to start transaction", "error", err)
		return nil, apperrors.WrapInternal(err, "failed to start purchase transaction")
	}
	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("failed to rollback transaction", "error", err)
		}
	}(tx) // Will be ignored if tx.Commit() succeeds

	// Step 1: Get and lock the skin
	skins, err := s.skinRepo.GetSkinsForUpdate(ctx, tx, []uuid.UUID{skinID})
	if err != nil {
		s.logger.Error("failed to get skin for update", "error", err, "skin_id", skinID)
		return nil, apperrors.WrapInternal(err, "failed to check skin availability")
	}
	if len(skins) == 0 {
		s.logger.Warn("skin not available for purchase", "skin_id", skinID)
		return nil, apperrors.NewNotFoundError("skin not available or already sold")
	}

	skinToPurchase := skins[0]
	s.logger.Info("skin locked for purchase", "skin_id", skinID, "price", skinToPurchase.Price)

	// Step 2: Get and check user's balance
	usr, err := s.userRepo.GetUserByIdForUpdate(ctx, tx, userID)
	if err != nil {
		s.logger.Error("failed to get user for update", "error", err, "user_id", userID)
		return nil, apperrors.WrapInternal(err, "failed to check user balance")
	}
	if usr.Balance < skinToPurchase.Price {
		s.logger.Warn("insufficient funds", "user_id", userID, "balance", usr.Balance, "required", skinToPurchase.Price)
		return nil, apperrors.NewValidationError("insufficient funds")
	}

	// Step 3: Create order record
	now := time.Now()
	ord := &models.Order{
		ID:          uuid.New(),
		UserID:      userID,
		TotalAmount: skinToPurchase.Price,
		Status:      models.OrderStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.orderRepo.Create(ctx, tx, ord); err != nil {
		s.logger.Error("failed to create order", "error", err, "order_id", ord.ID)
		return nil, apperrors.WrapInternal(err, "failed to create order")
	}

	// Step 4: Create order item
	orderItem := &models.OrderItem{
		ID:        uuid.New(),
		OrderID:   ord.ID,
		SkinID:    skinID,
		Price:     skinToPurchase.Price,
		CreatedAt: now,
	}

	if err := s.orderRepo.CreateOrderItem(ctx, tx, orderItem); err != nil {
		s.logger.Error("failed to create order item", "error", err, "order_id", ord.ID)
		return nil, apperrors.WrapInternal(err, "failed to create order item")
	}

	// Step 5: Update user balance
	newBalance := usr.Balance - skinToPurchase.Price
	if err := s.userRepo.UpdateBalance(ctx, tx, userID, newBalance); err != nil {
		s.logger.Error("failed to update user balance", "error", err, "user_id", userID)
		return nil, apperrors.WrapInternal(err, "failed to update balance")
	}

	// Step 6: Transfer skin ownership
	if err := s.skinRepo.UpdateOwnership(ctx, tx, []uuid.UUID{skinID}, userID); err != nil {
		s.logger.Error("failed to update skin ownership", "error", err, "skin_id", skinID)
		return nil, apperrors.WrapInternal(err, "failed to transfer skin ownership")
	}

	// Step 7: Update order status to completed
	ord.Status = models.OrderStatusCompleted
	if err := s.orderRepo.UpdateStatus(ctx, tx, ord.ID, models.OrderStatusCompleted); err != nil {
		s.logger.Error("failed to update order status", "error", err, "order_id", ord.ID)
		return nil, apperrors.WrapInternal(err, "failed to complete order")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", "error", err, "order_id", ord.ID)
		return nil, apperrors.WrapInternal(err, "failed to complete purchase")
	}

	s.logger.Info("skin purchase completed successfully",
		"user_id", userID,
		"skin_id", skinID,
		"order_id", ord.ID,
		"amount", skinToPurchase.Price)

	return ord, nil
}

// ListAvailableSkins returns all skins that can be purchased
func (s *MarketplaceService) ListAvailableSkins(ctx context.Context) ([]*models.Skin, error) {
	skins, err := s.skinRepo.GetAvailableSkins(ctx)
	if err != nil {
		return nil, apperrors.WrapInternal(err, "failed to list available skins")
	}
	return skins, nil
}

// ListUserSkins returns all skins owned by the given user
func (s *MarketplaceService) ListUserSkins(ctx context.Context, userID uuid.UUID) ([]*models.Skin, error) {
	skins, err := s.skinRepo.GetUserSkins(ctx, userID)
	if err != nil {
		return nil, apperrors.WrapInternal(err, "failed to list user skins")
	}
	return skins, nil
}

// GetOrder returns an order by id
func (s *MarketplaceService) GetOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	ord, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, apperrors.WrapInternal(err, "failed to get order")
	}
	if ord == nil {
		return nil, apperrors.NewNotFoundError("order not found")
	}
	return ord, nil
}
