package handlers

import (
	"net/http"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *services.User
}

func NewUserHandler(svc *services.User) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req models.UserSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, apperrors.NewValidationError(err.Error()))
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.svc.CreateUser(c.Request.Context(), user); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"email": user.Email})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	token, err := h.svc.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
