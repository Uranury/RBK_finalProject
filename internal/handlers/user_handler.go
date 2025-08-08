package handlers

import (
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	svc *services.User
}

func NewUserHandler(svc *services.User) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, err)
		return
	}

	user := &models.User{
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
	var req models.UserRequest
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
