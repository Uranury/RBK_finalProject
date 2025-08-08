package handlers

import (
	"errors"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"` // Optional: include app error code
}

// HandleError converts application errors to appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		httpStatus := mapErrorCodeToHTTPStatus(appErr.Code)

		c.JSON(httpStatus, ErrorResponse{
			Error: appErr.Message,
			Code:  int(appErr.Code), // Optional: include app-specific error code
		})
		return
	}

	// Handle unknown errors
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "Internal server error",
	})
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
