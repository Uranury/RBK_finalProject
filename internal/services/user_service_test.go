package services

import (
	"context"
	"io"
	"testing"

	"log/slog"

	"github.com/Uranury/RBK_finalProject/internal/auth"
	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of user.Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockUserRepository) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserRepository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, newBalance float64) error {
	args := m.Called(ctx, tx, userID, newBalance)
	return args.Error(0)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByIdForUpdate(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, tx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	// Create a logger that writes to io.Discard to avoid nil pointer issues
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService := auth.NewService("test-secret")

	tests := []struct {
		name          string
		user          *models.User
		mockSetup     func(*MockUserRepository)
		expectedError error
	}{
		{
			name: "successful user creation",
			user: &models.User{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockSetup: func(repo *MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "user already exists",
			user: &models.User{
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Existing User",
			},
			mockSetup: func(repo *MockUserRepository) {
				existingUser := &models.User{
					ID:    uuid.New(),
					Email: "existing@example.com",
					Name:  "Existing User",
				}
				repo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectedError: apperrors.ErrUserExists,
		},
		{
			name: "repository error on find by email",
			user: &models.User{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockSetup: func(repo *MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, assert.AnError)
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			service := NewUser(mockRepo, authService, logger)
			err := service.CreateUser(context.Background(), tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				// For wrapped errors, we check if the original error is contained
				if tt.expectedError == assert.AnError {
					assert.Contains(t, err.Error(), "failed to check existing user")
				} else {
					assert.Equal(t, tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, tt.user.ID)
				assert.Equal(t, 0.0, tt.user.Balance)
				assert.Equal(t, auth.User, tt.user.Role)
				assert.NotEmpty(t, tt.user.Password) // Should be hashed
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	// Create a logger that writes to io.Discard to avoid nil pointer issues
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService := auth.NewService("test-secret")

	// Create a properly hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Create a test user with properly hashed password
	testUser := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
		Role:     auth.User,
	}

	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*MockUserRepository)
		expectedError error
		expectToken   bool
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(repo *MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, "test@example.com").Return(testUser, nil)
			},
			expectedError: nil,
			expectToken:   true,
		},
		{
			name:     "user not found",
			email:    "nonexistent@example.com",
			password: "password123",
			mockSetup: func(repo *MockUserRepository) {
				repo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, nil)
			},
			expectedError: apperrors.ErrUserNotFound,
			expectToken:   false,
		},
		{
			name:     "empty email",
			email:    "",
			password: "password123",
			mockSetup: func(repo *MockUserRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("email and password is required"),
			expectToken:   false,
		},
		{
			name:     "empty password",
			email:    "test@example.com",
			password: "",
			mockSetup: func(repo *MockUserRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("email and password is required"),
			expectToken:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			service := NewUser(mockRepo, authService, logger)
			token, err := service.LoginUser(context.Background(), tt.email, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.expectToken {
					assert.NotEmpty(t, token)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserProfile(t *testing.T) {
	// Create a logger that writes to io.Discard to avoid nil pointer issues
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	authService := auth.NewService("test-secret")

	testUserID := uuid.New()
	testProfile := &models.UserProfile{
		Email:   "test@example.com",
		Name:    "Test User",
		Balance: 100.0,
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(*MockUserRepository)
		expectedError error
		expectedUser  *models.UserProfile
	}{
		{
			name:   "successful profile retrieval",
			userID: testUserID,
			mockSetup: func(repo *MockUserRepository) {
				repo.On("GetUserProfile", mock.Anything, testUserID).Return(testProfile, nil)
			},
			expectedError: nil,
			expectedUser:  testProfile,
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(repo *MockUserRepository) {
				repo.On("GetUserProfile", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil)
			},
			expectedError: apperrors.ErrUserNotFound,
			expectedUser:  nil,
		},
		{
			name:   "repository error",
			userID: testUserID,
			mockSetup: func(repo *MockUserRepository) {
				repo.On("GetUserProfile", mock.Anything, testUserID).Return(nil, assert.AnError)
			},
			expectedError: assert.AnError,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			service := NewUser(mockRepo, authService, logger)
			user, err := service.GetUserProfile(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				// For wrapped errors, we check if the original error is contained
				if tt.expectedError == assert.AnError {
					assert.Contains(t, err.Error(), "failed to fetch user profile")
				} else {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
