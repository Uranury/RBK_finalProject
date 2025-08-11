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

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account with email, password, and name
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserSignupRequest true "User registration data"
// @Success 201 {object} map[string]string "User created successfully"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /signup [post]
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

// Login godoc
// @Summary Authenticate user
// @Description Login with email and password to receive JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} map[string]string "Login successful"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /login [post]
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
