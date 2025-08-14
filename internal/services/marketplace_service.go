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

// TODO: Add asynq background workers later

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

func (s *MarketplaceService) PurchaseSkin(ctx context.Context, userID uuid.UUID, skinID uuid.UUID) (*models.Order, error) {
	s.logger.Info("starting skin purchase", "user_id", userID, "skin_id", skinID)

	// Start database transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return nil, apperrors.WrapInternal(err, "failed to begin transaction")
	}
	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("failed to rollback transaction", "error", err)
		}
	}(tx)

	// Step 1: Get and lock the skin
	skins, err := s.skinRepo.GetSkinsForUpdate(ctx, tx, []uuid.UUID{skinID})
	if err != nil {
		s.logger.Error("failed to get skin for update", "error", err, "skin_id", skinID)
		return nil, apperrors.WrapInternal(err, "failed to get skin for update")
	}
	if len(skins) == 0 {
		s.logger.Warn("skin not found for purchase", "skin_id", skinID)
		return nil, apperrors.NewNotFoundError("skin not found")
	}

	skinToPurchase := skins[0]

	// Check if skin is available for purchase
	if !skinToPurchase.Available {
		s.logger.Warn("skin is not available for purchase", "skin_id", skinID)
		return nil, apperrors.NewValidationError("skin is not available for purchase")
	}

	// Prevent users from buying their own skins
	if skinToPurchase.OwnerID != nil && *skinToPurchase.OwnerID == userID {
		s.logger.Warn("user attempted to buy their own skin", "user_id", userID, "skin_id", skinID)
		return nil, apperrors.NewValidationError("cannot purchase your own skin")
	}

	s.logger.Info("skin locked for purchase", "skin_id", skinID, "price", skinToPurchase.Price, "has_owner", skinToPurchase.OwnerID != nil)

	// Step 2: Get and check buyer's balance
	buyer, err := s.userRepo.GetUserByIdForUpdate(ctx, tx, userID)
	if err != nil {
		s.logger.Error("failed to get buyer for update", "error", err, "user_id", userID)
		return nil, apperrors.WrapInternal(err, "failed to get buyer for update")
	}
	if buyer.Balance < skinToPurchase.Price {
		s.logger.Warn("insufficient funds", "user_id", userID, "balance", buyer.Balance, "required", skinToPurchase.Price)
		return nil, apperrors.NewValidationError("insufficient funds")
	}

	// Step 3: If skin has an owner, get and lock the owner for balance update
	var owner *models.User
	if skinToPurchase.OwnerID != nil {
		owner, err = s.userRepo.GetUserByIdForUpdate(ctx, tx, *skinToPurchase.OwnerID)
		if err != nil {
			s.logger.Error("failed to get owner for update", "error", err, "owner_id", *skinToPurchase.OwnerID)
			return nil, apperrors.WrapInternal(err, "failed to get owner for update")
		}
		s.logger.Info("skin owner locked for payment", "owner_id", owner.ID, "current_balance", owner.Balance)
	}

	// Step 4: Create order record
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

	// Step 5: Create order item
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

	// Step 6: Update buyer balance (deduct payment)
	newBuyerBalance := buyer.Balance - skinToPurchase.Price
	if err := s.userRepo.UpdateBalance(ctx, tx, userID, newBuyerBalance); err != nil {
		s.logger.Error("failed to update buyer balance", "error", err, "user_id", userID)
		return nil, apperrors.WrapInternal(err, "failed to update buyer balance")
	}
	s.logger.Info("buyer balance updated", "user_id", userID, "old_balance", buyer.Balance, "new_balance", newBuyerBalance)

	// Step 7: If skin has an owner, credit them with the sale amount
	if owner != nil {
		newOwnerBalance := owner.Balance + skinToPurchase.Price
		if err := s.userRepo.UpdateBalance(ctx, tx, owner.ID, newOwnerBalance); err != nil {
			s.logger.Error("failed to update owner balance", "error", err, "owner_id", owner.ID)
			return nil, apperrors.WrapInternal(err, "failed to update owner balance")
		}
		s.logger.Info("owner credited with sale", "owner_id", owner.ID, "old_balance", owner.Balance, "new_balance", newOwnerBalance, "amount", skinToPurchase.Price)
	} else {
		s.logger.Info("no owner to credit - market-created skin", "skin_id", skinID)
	}

	// Step 8: Mark skin as not available (remove from listing) and transfer ownership
	if err := s.skinRepo.UpdateAvailability(ctx, tx, skinID, false); err != nil {
		s.logger.Error("failed to update skin availability", "error", err, "skin_id", skinID)
		return nil, apperrors.WrapInternal(err, "failed to update skin availability")
	}

	if err := s.skinRepo.UpdateOwnership(ctx, tx, []uuid.UUID{skinID}, userID); err != nil {
		s.logger.Error("failed to update skin ownership", "error", err, "skin_id", skinID)
		return nil, apperrors.WrapInternal(err, "failed to update skin ownership")
	}

	// Step 9: Update order status to completed
	ord.Status = models.OrderStatusCompleted
	if err := s.orderRepo.UpdateStatus(ctx, tx, ord.ID, models.OrderStatusCompleted); err != nil {
		s.logger.Error("failed to update order status", "error", err, "order_id", ord.ID)
		return nil, apperrors.WrapInternal(err, "failed to update order status")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", "error", err, "order_id", ord.ID)
		return nil, apperrors.WrapInternal(err, "failed to commit transaction")
	}

	s.logger.Info("skin purchase completed successfully",
		"user_id", userID,
		"skin_id", skinID,
		"order_id", ord.ID,
		"amount", skinToPurchase.Price,
		"owner_credited", owner != nil)

	return ord, nil
}

func (s *MarketplaceService) RemoveSkinFromListing(ctx context.Context, userID uuid.UUID, skinID uuid.UUID) error {
	s.logger.Info("starting skin removal from listing", "user_id", userID, "skin_id", skinID)

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return apperrors.WrapInternal(err, "failed to begin transaction")
	}

	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("failed to rollback transaction", "error", err)
		}
	}(tx)

	skins, err := s.skinRepo.GetSkinsForUpdate(ctx, tx, []uuid.UUID{skinID})
	if err != nil {
		s.logger.Error("failed to get skin for update", "error", err, "skin_id", skinID)
		return apperrors.WrapInternal(err, "failed to get skin for update")
	}

	if len(skins) == 0 {
		s.logger.Warn("skin not found", "skin_id", skinID)
		return apperrors.NewNotFoundError("skin not found")
	}

	skinToRemove := skins[0]
	if skinToRemove.OwnerID == nil || *skinToRemove.OwnerID != userID {
		s.logger.Warn("user doesn't own this skin", "user_id", userID, "skin_id", skinID, "actual_owner", skinToRemove.OwnerID)
		return apperrors.NewForbiddenError("you can only remove your own skins from listing")
	}

	if !skinToRemove.Available {
		s.logger.Warn("skin is not currently listed", "skin_id", skinID)
		return apperrors.NewValidationError("skin is not currently listed for sale")
	}

	if err := s.skinRepo.UpdateAvailability(ctx, tx, skinToRemove.ID, false); err != nil {
		s.logger.Error("failed to update skin availability", "error", err, "skin_id", skinID)
		return apperrors.WrapInternal(err, "failed to update skin availability")
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", "error", err, "skin_id", skinID)
		return apperrors.WrapInternal(err, "failed to commit transaction")
	}

	s.logger.Info("skin removed from listing successfully", "user_id", userID, "skin_id", skinID)
	return nil
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
