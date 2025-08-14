package apperrors

import (
	"errors"
	"fmt"
)

// ErrorCode represents application-specific error codes
type ErrorCode int

const (
	CodeInternal ErrorCode = iota + 1000
	CodeNotFound
	CodeAlreadyExists
	CodeInvalidInput
	CodeUnauthorized
	CodeForbidden
	CodeValidation
)

// AppError represents an application error with code and context
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"` // Don't expose internal errors in JSON
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work
func (e *AppError) Unwrap() error {
	return e.Err
}

// Constructor functions for common errors

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: message,
	}
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    CodeForbidden,
		Message: message,
	}
}

func NewAlreadyExistsError(message string) *AppError {
	return &AppError{
		Code:    CodeAlreadyExists,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: message,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Message: message,
		Err:     err,
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, we might want to preserve or modify it
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapInternal Convenience function to wrap with internal error code
func WrapInternal(err error, message string) *AppError {
	return WrapError(err, CodeInternal, message)
}

// Predefined common errors
var (
	ErrUserNotFound       = NewNotFoundError("user not found")
	ErrUserExists         = NewAlreadyExistsError("user already exists")
	ErrInvalidCredentials = NewUnauthorizedError("invalid credentials")
)
