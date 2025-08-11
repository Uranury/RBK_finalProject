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
	SkinID string `json:"skin_id" binding:"required"`
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

func (h *MarketplaceHandler) ListAvailable(c *gin.Context) {
	skins, err := h.svc.ListAvailableSkins(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, skins)
}

func (h *MarketplaceHandler) ListMine(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	skins, err := h.svc.ListUserSkins(c.Request.Context(), userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, skins)
}

func (h *MarketplaceHandler) GetOrder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	ord, err := h.svc.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		HandleError(c, err)
		return
	}
	if ord.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, ord)
}
