package services

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/Uranury/RBK_finalProject/internal/auth"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/internal/repositories/user"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	repo   user.Repository
	Auth   auth.Service
	logger *slog.Logger
}

func NewUser(repo user.Repository, Auth *auth.Service, logger *slog.Logger) *User {
	return &User{repo: repo, Auth: *Auth, logger: logger}
}

func (s *User) CreateUser(ctx context.Context, user *models.User) error {
	existingUser, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		s.logger.Error("failed to check existing user", "email", user.Email, "error", err)
		return apperrors.WrapInternal(err, "failed to check existing user")
	}
	if existingUser != nil {
		s.logger.Warn("attempt to create duplicate user", "email", user.Email)
		return apperrors.ErrUserExists
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.ID = uuid.New()
	user.Balance = 0.0
	user.Role = auth.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", "user_id", user.ID, "error", err)
		return apperrors.NewInternalError("failed to hash password", err)
	}
	user.Password = string(hashedPassword)

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("failed to create user in repository", "user_id", user.ID, "error", err)
		return apperrors.WrapInternal(err, "failed to create user")
	}

	s.logger.Info("user created successfully", "user_id", user.ID, "email", user.Email)
	return nil
}

func (s *User) LoginUser(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		s.logger.Warn("login attempt with missing email or password")
		return "", apperrors.NewValidationError("email and password is required")
	}

	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to find user by email", "email", email, "error", err)
		return "", err
	}
	if existingUser == nil {
		s.logger.Warn("login attempt for non-existent user", "email", email)
		return "", apperrors.ErrUserNotFound
	}

	password = strings.TrimSpace(password)
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		s.logger.Warn("invalid login credentials", "email", email)
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := s.Auth.GenerateJWT(existingUser.ID, existingUser.Role)
	if err != nil {
		s.logger.Error("failed to generate JWT", "user_id", existingUser.ID, "error", err)
		return "", apperrors.WrapInternal(err, "failed to generate JWT")
	}

	s.logger.Info("user logged in successfully", "user_id", existingUser.ID, "email", email)
	return token, nil
}

func (s *User) GetUserProfile(ctx context.Context, id uuid.UUID) (*models.UserProfile, error) {
	s.logger.Info("fetching user profile", "user_id", id)

	usr, err := s.repo.GetUserProfile(ctx, id)
	if err != nil {
		s.logger.Error("failed to fetch user profile", "user_id", id, "error", err)
		return nil, apperrors.WrapInternal(err, "failed to fetch user profile")
	}
	if usr == nil {
		s.logger.Info("user profile not found", "user_id", id)
		return nil, apperrors.ErrUserNotFound
	}

	return usr, nil
}
