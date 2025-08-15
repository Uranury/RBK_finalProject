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

// Purchase godoc
// @Summary Purchase a skin
// @Description Purchase a skin from the marketplace using user's balance
// @Tags marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param purchase body purchaseRequest true "Purchase request"
// @Success 201 {object} models.Order "Purchase successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Skin not available"
// @Failure 422 {object} ErrorResponse "Insufficient funds"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /marketplace/purchase [post]
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

type sellRequest struct {
	SkinID string  `json:"skin_id" binding:"required"`
	Price  float64 `json:"price" binding:"required,gt=0,lte=1000000" description:"price must be > 0 and <= 1,000,000"`
}

// Sell godoc
// @Summary Sell a skin
// @Description Sell a skin that you own. Price must be > 0 and <= 1,000,000.
// @Tags marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param purchase body sellRequest true "Sell request"
// @Success 201 {string} string "UUID of listed skin"
// @Failure 400 {object} ErrorResponse "Invalid request (e.g., invalid skinID, invalid price, skin already listed)"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden: skin ownership required"
// @Failure 404 {object} ErrorResponse "Skin not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /marketplace/sell [post]
func (h *MarketplaceHandler) Sell(c *gin.Context) {
	var req sellRequest
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

	if err := h.svc.SellSkin(c.Request.Context(), userID, skinID, req.Price); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, skinID)
}

// RemoveFromListing godoc
// @Summary Remove a skin
// @Description Remove a user's skin from listing
// @Tags marketplace
// @Produce json
// @Param skin_id query string true "Skin ID to remove from listing"
// @Success 200 {string} string "Skin ID removed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Skin not available"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /marketplace/skins/remove [get]
func (h *MarketplaceHandler) RemoveFromListing(c *gin.Context) {
	querySkinID := c.Query("skin_id")
	if querySkinID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "skin_id is required"})
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skinID, err := uuid.Parse(querySkinID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid skinId"})
	}

	err = h.svc.RemoveSkinFromListing(c.Request.Context(), userID, skinID)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, skinID.String())
}

// ListAvailable godoc
// @Summary List available skins
// @Description Get all skins available for purchase in the marketplace
// @Tags marketplace
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Skin "List of available skins"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /marketplace/skins [get]
func (h *MarketplaceHandler) ListAvailable(c *gin.Context) {
	skins, err := h.svc.ListAvailableSkins(c.Request.Context())
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, skins)
}

// ListMine godoc
// @Summary List user's skins
// @Description Get all skins owned by the authenticated user
// @Tags marketplace
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Skin "List of user's skins"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /marketplace/skins/mine [get]
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

// GetOrder godoc
// @Summary Get order details
// @Description Get details of a specific order by ID (only if owned by the user)
// @Tags marketplace
// @Produce json
// @Security BearerAuth
// @Param order_id path string true "Order ID" format(uuid)
// @Success 200 {object} models.Order "Order details"
// @Failure 400 {object} ErrorResponse "Invalid order ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - order not owned by user"
// @Failure 404 {object} ErrorResponse "Order not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /marketplace/orders/{order_id} [get]
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
