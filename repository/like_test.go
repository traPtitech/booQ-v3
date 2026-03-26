package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

func TestLikeRepository_GetByItemID(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) (int, []*domain.Like)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) (int, []*domain.Like) {
				models := []*like{
					{ItemID: 1, UserID: "user2"},
					{ItemID: 1, UserID: "user1"},
					{ItemID: 2, UserID: "user3"},
				}
				for _, model := range models {
					if err := db.Create(model).Error; err != nil {
						t.Fatalf("failed to create test like: %v", err)
					}
				}

				return 1, []*domain.Like{
					{ItemID: 1, UserID: "user1"},
					{ItemID: 1, UserID: "user2"},
				}
			},
			wantErr: false,
		},
		{
			name: "success: empty result",
			setup: func(t *testing.T, db *gorm.DB) (int, []*domain.Like) {
				return 9999, []*domain.Like{}
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewLikeRepository(db)
			itemID, expected := tc.setup(t, db)

			likes, err := repo.GetByItemID(itemID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.ElementsMatch(t, expected, likes)
		})
	}
}

func TestLikeRepository_Exists(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) (int, string)
		expected bool
		wantErr  bool
	}{
		{
			name: "success: exists",
			setup: func(t *testing.T, db *gorm.DB) (int, string) {
				model := &like{ItemID: 1, UserID: "user1"}
				db.Create(model)

				return 1, "user1"
			},
			expected: true,
			wantErr:  false,
		},
		{
			name: "success: not exists",
			setup: func(t *testing.T, db *gorm.DB) (int, string) {
				return 1, "unknown"
			},
			expected: false,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewLikeRepository(db)
			itemID, userID := tc.setup(t, db)

			exists, err := repo.Exists(itemID, userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, exists)
		})
	}
}

func TestLikeRepository_Create(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB)
		like    *domain.Like
		wantErr bool
	}{
		{
			name: "success",
			like: &domain.Like{ItemID: 1, UserID: "user1"},
			setup: func(t *testing.T, db *gorm.DB) {
				// Setup code if needed
			},
			wantErr: false,
		},
		{
			name: "failure: duplicated like",
			like: &domain.Like{ItemID: 2, UserID: "user2"},
			setup: func(t *testing.T, db *gorm.DB) {
				err := db.Create(&like{ItemID: 2, UserID: "user2"}).Error
				assert.NoError(t, err)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewLikeRepository(db)
			tc.setup(t, db)

			err := repo.Create(tc.like)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			exists, err := repo.Exists(tc.like.ItemID, tc.like.UserID)
			assert.NoError(t, err)
			assert.True(t, exists)
		})
	}
}

func TestLikeRepository_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		setup     func(t *testing.T, db *gorm.DB) (int, string)
		expected  error
		wantError bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) (int, string) {
				model := &like{ItemID: 1, UserID: "user1"}
				db.Create(model)

				return 1, "user1"
			},
			expected: nil,
		},
		{
			name: "failure: not found",
			setup: func(t *testing.T, db *gorm.DB) (int, string) {
				return 9999, "user1"
			},
			expected:  domain.ErrNotFound,
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewLikeRepository(db)
			itemID, userID := tc.setup(t, db)

			err := repo.Delete(itemID, userID)
			if tc.wantError {
				assert.ErrorIs(t, err, tc.expected)
				return
			}
			assert.NoError(t, err)

			exists, err := repo.Exists(itemID, userID)
			assert.NoError(t, err)
			assert.False(t, exists)
		})
	}
}
