package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error   string            `json:"error"`
	Code    int               `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// HandleError converts application errors to appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	log.Printf("HandleError received error: %v (type: %T)", err, err)

	// Handle JSON parsing & binding errors first
	if errors.Is(err, io.EOF) || strings.Contains(err.Error(), "json:") {
		handleBindingError(c, err)
		return
	}

	// Handle struct validation errors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		handleValidationErrors(c, validationErrors)
		return
	}

	// Handle generic binding/validation errors
	if strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "binding") {
		handleBindingError(c, err)
		return
	}

	// Handle application-specific errors
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		httpStatus := mapErrorCodeToHTTPStatus(appErr.Code)
		c.JSON(httpStatus, ErrorResponse{
			Error: appErr.Message,
			Code:  int(appErr.Code),
		})
		return
	}

	// Unknown errors — don’t leak internal details
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "An unexpected error occurred. Please try again later.",
	})
}

// handleValidationErrors handles Gin validator validation errors
func handleValidationErrors(c *gin.Context, validationErrors validator.ValidationErrors) {
	details := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := strings.ToLower(fieldError.Field())
		message := getValidationMessage(fieldError.Tag(), fieldError.Param())
		details[field] = message
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "Validation failed",
		Code:    int(apperrors.CodeValidation),
		Details: details,
	})
}

// handleBindingError handles Gin binding errors (like JSON parsing issues)
func handleBindingError(c *gin.Context, err error) {
	var message string

	switch {
	case strings.Contains(err.Error(), "json"):
		message = "Invalid JSON format"
	case strings.Contains(err.Error(), "binding"):
		message = "Invalid request data"
	default:
		message = "Invalid input data"
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: message,
		Code:  int(apperrors.CodeValidation),
	})
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Please provide a valid email address"
	case "min":
		return "Value must be at least " + param
	case "max":
		return "Value must be at most " + param
	case "gt":
		return "Value must be greater than " + param
	case "gte":
		return "Value must be greater than or equal to " + param
	case "lt":
		return "Value must be less than " + param
	case "lte":
		return "Value must be less than or equal to " + param
	case "oneof":
		return "Value must be one of: " + param
	case "numeric":
		return "Value must be a number"
	case "alpha":
		return "Value must contain only letters"
	case "alphanum":
		return "Value must contain only letters and numbers"
	default:
		return "Invalid value"
	}
}

// mapErrorCodeToHTTPStatus maps application error codes to HTTP status codes
func mapErrorCodeToHTTPStatus(code apperrors.ErrorCode) int {
	switch code {
	case apperrors.CodeNotFound:
		return http.StatusNotFound // 404
	case apperrors.CodeAlreadyExists:
		return http.StatusConflict // 409
	case apperrors.CodeInvalidInput, apperrors.CodeValidation:
		return http.StatusBadRequest // 400
	case apperrors.CodeUnauthorized:
		return http.StatusUnauthorized // 401
	case apperrors.CodeForbidden:
		return http.StatusForbidden // 403
	case apperrors.CodeInternal:
		return http.StatusInternalServerError // 500
	default:
		return http.StatusInternalServerError // 500
	}
}
