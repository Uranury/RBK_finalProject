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
	Gun       string  `json:"gun" binding:"required"`
	Rarity    string  `json:"rarity" binding:"required"`
	Condition float64 `json:"condition" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	Image     string  `json:"image"`
}

// GetGuns godoc
// @Summary Get all available guns
// @Description Get a list of all available guns in the system
// @Tags skins
// @Produce json
// @Success 200 {array} models.Gun "List of available guns"
// @Router /guns [get]
func (h *SkinHandler) GetGuns(c *gin.Context) {
	guns := h.svc.GetAllGuns()
	c.JSON(http.StatusOK, guns)
}

// Create godoc
// @Summary Create a new skin
// @Description Create a new skin and add it to the marketplace
// @Tags skins
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param skin body createSkinRequest true "Skin creation data"
// @Success 201 {object} models.Skin "Skin created successfully"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /skins [post]
func (h *SkinHandler) Create(c *gin.Context) {
	var req createSkinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	skin := &models.Skin{
		Name:      req.Name,
		Gun:       models.Gun(req.Gun),
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
