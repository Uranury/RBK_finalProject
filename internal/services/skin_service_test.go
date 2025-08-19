package services

import (
	"context"
	"io"
	"testing"

	"log/slog"

	"github.com/Uranury/RBK_finalProject/internal/models"
	"github.com/Uranury/RBK_finalProject/pkg/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSkinRepository is a mock implementation of skin.Repository
type MockSkinRepository struct {
	mock.Mock
}

func (m *MockSkinRepository) Create(ctx context.Context, skin *models.Skin) error {
	args := m.Called(ctx, skin)
	return args.Error(0)
}

func (m *MockSkinRepository) GetSkin(ctx context.Context, id uuid.UUID) (*models.Skin, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Skin), args.Error(1)
}

func (m *MockSkinRepository) GetSkinsForSellUpdate(ctx context.Context, tx *sqlx.Tx, ids []uuid.UUID) ([]*models.Skin, error) {
	args := m.Called(ctx, tx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Skin), args.Error(1)
}

func (m *MockSkinRepository) UpdatePrice(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, price float64) error {
	args := m.Called(ctx, tx, id, price)
	return args.Error(0)
}

func (m *MockSkinRepository) GetAvailableSkins(ctx context.Context) ([]*models.Skin, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Skin), args.Error(1)
}

func (m *MockSkinRepository) GetUserSkins(ctx context.Context, userID uuid.UUID) ([]*models.Skin, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Skin), args.Error(1)
}

// Add missing methods to satisfy the interface
func (m *MockSkinRepository) GetSkinsForUpdate(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID) ([]*models.Skin, error) {
	args := m.Called(ctx, tx, skinIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Skin), args.Error(1)
}

func (m *MockSkinRepository) UpdateOwnership(ctx context.Context, tx *sqlx.Tx, skinIDs []uuid.UUID, newOwnerID uuid.UUID) error {
	args := m.Called(ctx, tx, skinIDs, newOwnerID)
	return args.Error(0)
}

func (m *MockSkinRepository) UpdateForSale(ctx context.Context, tx *sqlx.Tx, skinID uuid.UUID, price float64, available bool) error {
	args := m.Called(ctx, tx, skinID, price, available)
	return args.Error(0)
}

func (m *MockSkinRepository) UpdateAvailability(ctx context.Context, tx *sqlx.Tx, skinID uuid.UUID, available bool) error {
	args := m.Called(ctx, tx, skinID, available)
	return args.Error(0)
}

func TestSkinService_GetAllGuns(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockRepo := new(MockSkinRepository)
	service := NewSkin(mockRepo, logger)

	guns := service.GetAllGuns()

	// Test that we get a non-empty list
	assert.NotEmpty(t, guns)

	// Test that specific guns are included
	gunMap := make(map[models.Gun]bool)
	for _, gun := range guns {
		gunMap[gun] = true
	}

	// Check for some expected guns
	assert.True(t, gunMap[models.AK47])
	assert.True(t, gunMap[models.AWP])
	assert.True(t, gunMap[models.Karambit])
	assert.True(t, gunMap[models.DesertEagle])

	// Test that we have a reasonable number of guns
	assert.GreaterOrEqual(t, len(guns), 40) // Should have many guns
}

func TestSkinService_GetAllWears(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mockRepo := new(MockSkinRepository)
	service := NewSkin(mockRepo, logger)

	wears := service.GetAllWears()

	// Test that we get exactly 5 wear levels
	assert.Len(t, wears, 5)

	// Test that all expected wear levels are present
	wearMap := make(map[models.Wear]bool)
	for _, wear := range wears {
		wearMap[wear] = true
	}

	assert.True(t, wearMap[models.FactoryNew])
	assert.True(t, wearMap[models.MinimalWear])
	assert.True(t, wearMap[models.FieldTested])
	assert.True(t, wearMap[models.WellWorn])
	assert.True(t, wearMap[models.BattleScarred])
}

func TestSkinService_CreateSkin(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name          string
		skin          *models.Skin
		mockSetup     func(*MockSkinRepository)
		expectedError error
		checkSkin     func(*testing.T, *models.Skin)
	}{
		{
			name: "successful skin creation",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     100.0,
				Condition: 0.8,
				Gun:       models.AK47,
			},
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Skin")).Return(nil)
			},
			expectedError: nil,
			checkSkin: func(t *testing.T, skin *models.Skin) {
				assert.NotEqual(t, uuid.Nil, skin.ID)
				assert.True(t, skin.Available)
				assert.Nil(t, skin.OwnerID)
				assert.Equal(t, models.BattleScarred, skin.Wear) // 0.8 condition should be Battle-Scarred
			},
		},
		{
			name: "successful skin creation with minimal wear",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     100.0,
				Condition: 0.1, // This should be Minimal Wear
				Gun:       models.AK47,
			},
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Skin")).Return(nil)
			},
			expectedError: nil,
			checkSkin: func(t *testing.T, skin *models.Skin) {
				assert.NotEqual(t, uuid.Nil, skin.ID)
				assert.True(t, skin.Available)
				assert.Nil(t, skin.OwnerID)
				assert.Equal(t, models.MinimalWear, skin.Wear) // 0.1 condition should be Minimal Wear
			},
		},
		{
			name: "missing name",
			skin: &models.Skin{
				Rarity:    "Classified",
				Price:     100.0,
				Condition: 0.8,
			},
			mockSetup: func(repo *MockSkinRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("skin name is required"),
			checkSkin:     nil,
		},
		{
			name: "missing rarity",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Price:     100.0,
				Condition: 0.8,
			},
			mockSetup: func(repo *MockSkinRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("skin rarity is required"),
			checkSkin:     nil,
		},
		{
			name: "invalid price (zero)",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     0.0,
				Condition: 0.8,
			},
			mockSetup: func(repo *MockSkinRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("skin price can't be negative or zero"),
			checkSkin:     nil,
		},
		{
			name: "invalid price (negative)",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     -10.0,
				Condition: 0.8,
			},
			mockSetup: func(repo *MockSkinRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("skin price can't be negative or zero"),
			checkSkin:     nil,
		},
		{
			name: "invalid condition (too high)",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     100.0,
				Condition: 1.5,
			},
			mockSetup: func(repo *MockSkinRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("Skin condition must be between 0 and 1"),
			checkSkin:     nil,
		},
		{
			name: "invalid condition (negative)",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     100.0,
				Condition: -0.1,
			},
			mockSetup: func(repo *MockSkinRepository) {
				// No mock setup needed
			},
			expectedError: apperrors.NewValidationError("Skin condition must be between 0 and 1"),
			checkSkin:     nil,
		},
		{
			name: "repository error",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     100.0,
				Condition: 0.8,
			},
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Skin")).Return(assert.AnError)
			},
			expectedError: assert.AnError,
			checkSkin:     nil,
		},
		{
			name: "auto-set gun when not provided",
			skin: &models.Skin{
				Name:      "AK-47 | Redline",
				Rarity:    "Classified",
				Price:     100.0,
				Condition: 0.8,
				// Gun not set
			},
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Skin")).Return(nil)
			},
			expectedError: nil,
			checkSkin: func(t *testing.T, skin *models.Skin) {
				assert.Equal(t, models.AK47, skin.Gun) // Should default to AK47
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSkinRepository)
			tt.mockSetup(mockRepo)

			service := NewSkin(mockRepo, logger)
			result, err := service.CreateSkin(context.Background(), tt.skin)

			if tt.expectedError != nil {
				assert.Error(t, err)
				// For wrapped errors, we check if the original error is contained
				if tt.expectedError == assert.AnError {
					assert.Contains(t, err.Error(), "Failed to create skin")
				} else {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.checkSkin != nil {
					tt.checkSkin(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSkinService_GetSkinByID(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	testSkinID := uuid.New()
	testSkin := &models.Skin{
		ID:        testSkinID,
		Name:      "AK-47 | Redline",
		Rarity:    "Classified",
		Price:     100.0,
		Condition: 0.8,
		Gun:       models.AK47,
		Wear:      models.BattleScarred,
	}

	tests := []struct {
		name          string
		skinID        uuid.UUID
		mockSetup     func(*MockSkinRepository)
		expectedError error
		expectedSkin  *models.Skin
	}{
		{
			name:   "successful skin retrieval",
			skinID: testSkinID,
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("GetSkin", mock.Anything, testSkinID).Return(testSkin, nil)
			},
			expectedError: nil,
			expectedSkin:  testSkin,
		},
		{
			name:   "skin not found",
			skinID: uuid.New(),
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("GetSkin", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil)
			},
			expectedError: apperrors.NewNotFoundError("Skin not found"),
			expectedSkin:  nil,
		},
		{
			name:   "repository error",
			skinID: testSkinID,
			mockSetup: func(repo *MockSkinRepository) {
				repo.On("GetSkin", mock.Anything, testSkinID).Return(nil, assert.AnError)
			},
			expectedError: assert.AnError,
			expectedSkin:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSkinRepository)
			tt.mockSetup(mockRepo)

			service := NewSkin(mockRepo, logger)
			skin, err := service.GetSkinByID(context.Background(), tt.skinID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				// For wrapped errors, we check if the original error is contained
				if tt.expectedError == assert.AnError {
					assert.Contains(t, err.Error(), "Failed to retrieve skin")
				} else {
					assert.Equal(t, tt.expectedError, err)
				}
				assert.Nil(t, skin)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSkin, skin)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
