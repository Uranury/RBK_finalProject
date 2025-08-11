package handlers

import (
	"net/http"

	"github.com/Uranury/RBK_finalProject/internal/middleware"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MarketplaceHandler struct {
	svc *services.MarketplaceService
}

func NewMarketplaceHandler(svc *services.MarketplaceService) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc}
}

type purchaseRequest struct {
	SkinID string `json:"skinId" binding:"required"`
}

func (h *MarketplaceHandler) Purchase(c *gin.Context) {
	var req purchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skinID, err := uuid.Parse(req.SkinID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skinId"})
		return
	}

	order, err := h.svc.PurchaseSkin(c.Request.Context(), userID, skinID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, order)
}
