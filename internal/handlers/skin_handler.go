package handlers

import (
	"net/http"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/gin-gonic/gin"
)

type SkinHandler struct {
	svc *services.Skin
}

func NewSkinHandler(svc *services.Skin) *SkinHandler {
	return &SkinHandler{svc: svc}
}

type createSkinRequest struct {
	Name      string  `json:"name" binding:"required"`
	Rarity    string  `json:"rarity" binding:"required"`
	Condition float64 `json:"condition" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Image     string  `json:"image"`
}

func (h *SkinHandler) Create(c *gin.Context) {
	var req createSkinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	skin := &models.Skin{
		Name:      req.Name,
		Rarity:    req.Rarity,
		Condition: req.Condition,
		Price:     req.Price,
		Image:     req.Image,
	}

	created, err := h.svc.CreateSkin(c.Request.Context(), skin)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}
