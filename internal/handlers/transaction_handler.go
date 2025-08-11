package handlers

import (
	"net/http"

	"github.com/Uranury/RBK_finalProject/internal/middleware"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	svc *services.TransactionService
}

func NewTransactionHandler(svc *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

// Withdraw godoc
// @Summary Withdraw money from balance
// @Description Withdraw a specified amount from the user's balance
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param withdrawal body models.WithdrawRequest true "Withdrawal request"
// @Success 200 {object} models.Transaction "Withdrawal successful"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 422 {object} ErrorResponse "Insufficient funds"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transactions/withdraw [post]
func (h *TransactionHandler) Withdraw(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	transaction, err := h.svc.Withdraw(c.Request.Context(), userID, req.Amount)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Deposit godoc
// @Summary Deposit money to balance
// @Description Deposit a specified amount to the user's balance
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param deposit body models.DepositRequest true "Deposit request"
// @Success 200 {object} models.Transaction "Deposit successful"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transactions/deposit [post]
func (h *TransactionHandler) Deposit(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	transaction, err := h.svc.Deposit(c.Request.Context(), userID, req.Amount)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// GetHistory godoc
// @Summary Get transaction history
// @Description Get the transaction history for the authenticated user
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Transaction "Transaction history"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transactions/history [get]
func (h *TransactionHandler) GetHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transactions, err := h.svc.GetUserTransactions(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, transactions)
}
