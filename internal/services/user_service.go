package services

import (
	"context"
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
	repo user.Repository
	Auth auth.Service
}

func NewUser(repo user.Repository, Auth *auth.Service) *User {
	return &User{repo: repo, Auth: *Auth}
}

func (s *User) CreateUser(ctx context.Context, user *models.User) error {
	existingUser, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return apperrors.WrapInternal(err, "failed to check existing user")
	}
	if existingUser != nil {
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
		return apperrors.NewInternalError("failed to hash password", err)
	}
	user.Password = string(hashedPassword)

	if err := s.repo.Create(ctx, user); err != nil {
		return apperrors.WrapInternal(err, "failed to create user")
	}
	return nil
}

func (s *User) LoginUser(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", apperrors.NewValidationError("email and password is required")
	}

	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if existingUser == nil {
		return "", apperrors.NewNotFoundError("user not found")
	}

	password = strings.TrimSpace(password)
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	token, err := s.Auth.GenerateJWT(existingUser.ID, existingUser.Role)
	if err != nil {
		return "", apperrors.WrapInternal(err, "failed to generate JWT")
	}

	return token, nil
}

func (s *User) GetUserProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	//existingUser, err := s.repo.FindByID(ctx, id)
	//if err != nil {
	//	return nil, apperrors.WrapInternal(err, "failed to find user")
	//}
	//if existingUser == nil {
	//	return nil, apperrors.ErrUserNotFound
	//}
	return nil, nil
}
