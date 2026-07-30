package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

func TestOwnershipRepository_GetByID(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) (int, *domain.Ownership)
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) (int, *domain.Ownership) {
				createTestItems(t, db, 1)
				model := &ownership{ItemID: 1, UserID: "user1", Rentable: true, Memo: "memo"}
				if err := db.Create(model).Error; err != nil {
					t.Fatalf("failed to create test ownership: %v", err)
				}
				expected := &domain.Ownership{
					ID:       model.ID,
					ItemID:   1,
					UserID:   "user1",
					Rentable: true,
					Memo:     "memo",
				}
				return model.ID, expected
			},
			wantErr: nil,
		},
		{
			name: "failure: ownership not found",
			setup: func(t *testing.T, db *gorm.DB) (int, *domain.Ownership) {
				return 9999, nil
			},
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewOwnershipRepository(db)
			id, expected := tc.setup(t, db)

			ownership, err := repo.GetByID(id)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, ownership)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expected, ownership)
			}
		})
	}
}

func TestOwnershipRepository_GetByItemID(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) (int, []*domain.Ownership)
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) (int, []*domain.Ownership) {
				createTestItems(t, db, 1, 2)
				models := []*ownership{
					{ItemID: 1, UserID: "user1", Rentable: true, Memo: "memo1"},
					{ItemID: 1, UserID: "user2", Rentable: false, Memo: "memo2"},
					{ItemID: 2, UserID: "user3", Rentable: true, Memo: "memo3"},
				}
				for _, model := range models {
					if err := db.Create(model).Error; err != nil {
						t.Fatalf("failed to create test ownership: %v", err)
					}
				}
				expected := []*domain.Ownership{
					{ID: models[0].ID, ItemID: 1, UserID: "user1", Rentable: true, Memo: "memo1"},
					{ID: models[1].ID, ItemID: 1, UserID: "user2", Rentable: false, Memo: "memo2"},
				}
				return 1, expected
			},
			wantErr: nil,
		},
		{
			name: "success: empty result",
			setup: func(t *testing.T, db *gorm.DB) (int, []*domain.Ownership) {
				return 9999, []*domain.Ownership{}
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewOwnershipRepository(db)
			itemID, expected := tc.setup(t, db)

			ownerships, err := repo.GetByItemID(itemID)
			assert.NoError(t, err)
			assert.Equal(t, expected, ownerships)
		})
	}
}

func TestOwnershipRepository_GetByUserID(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) (string, []*domain.Ownership)
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) (string, []*domain.Ownership) {
				createTestItems(t, db, 1, 2, 3)
				models := []*ownership{
					{ItemID: 1, UserID: "target-user", Rentable: true, Memo: "memo1"},
					{ItemID: 2, UserID: "another-user", Rentable: false, Memo: "memo2"},
					{ItemID: 3, UserID: "target-user", Rentable: true, Memo: "memo3"},
				}
				for _, model := range models {
					if err := db.Create(model).Error; err != nil {
						t.Fatalf("failed to create test ownership: %v", err)
					}
				}
				expected := []*domain.Ownership{
					{ID: models[0].ID, ItemID: 1, UserID: "target-user", Rentable: true, Memo: "memo1"},
					{ID: models[2].ID, ItemID: 3, UserID: "target-user", Rentable: true, Memo: "memo3"},
				}
				return "target-user", expected
			},
			wantErr: nil,
		},
		{
			name: "success: empty result",
			setup: func(t *testing.T, db *gorm.DB) (string, []*domain.Ownership) {
				return "non-existent-user", []*domain.Ownership{}
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewOwnershipRepository(db)
			userID, expected := tc.setup(t, db)

			ownerships, err := repo.GetByUserID(userID)
			assert.NoError(t, err)
			assert.Equal(t, expected, ownerships)
		})
	}
}

func TestOwnershipRepository_Create(t *testing.T) {
	testCases := []struct {
		name      string
		ownership *domain.Ownership
		wantErr   error
	}{
		{
			name: "success",
			ownership: &domain.Ownership{
				ItemID:   10,
				UserID:   "user1",
				Rentable: true,
				Memo:     "created memo",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			createTestItems(t, db, tc.ownership.ItemID)
			repo := NewOwnershipRepository(db)

			created, err := repo.Create(tc.ownership)
			assert.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				assert.NotNil(t, created)
				assert.NotZero(t, created.ID)
				assert.Equal(t, tc.ownership.ItemID, created.ItemID)
				assert.Equal(t, tc.ownership.UserID, created.UserID)
				assert.Equal(t, tc.ownership.Rentable, created.Rentable)
				assert.Equal(t, tc.ownership.Memo, created.Memo)

				var model ownership
				err = db.First(&model, created.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tc.ownership.ItemID, model.ItemID)
				assert.Equal(t, tc.ownership.UserID, model.UserID)
				assert.Equal(t, tc.ownership.Rentable, model.Rentable)
				assert.Equal(t, tc.ownership.Memo, model.Memo)
			}
		})
	}
}

func TestOwnershipRepository_Update(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB) *domain.Ownership
		wantErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) *domain.Ownership {
				createTestItems(t, db, 10, 11)
				model := &ownership{ItemID: 10, UserID: "user1", Rentable: true, Memo: "before"}
				if err := db.Create(model).Error; err != nil {
					t.Fatalf("failed to create test ownership: %v", err)
				}
				return &domain.Ownership{
					ID:       model.ID,
					ItemID:   11,
					UserID:   "user1",
					Rentable: false,
					Memo:     "after",
				}
			},
			wantErr: nil,
		},
		{
			name: "success: create new ownership if not exist",
			setup: func(t *testing.T, db *gorm.DB) *domain.Ownership {
				createTestItems(t, db, 12)
				return &domain.Ownership{
					ID:       9999,
					ItemID:   12,
					UserID:   "user2",
					Rentable: true,
					Memo:     "new ownership",
				}
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewOwnershipRepository(db)
			updateData := tc.setup(t, db)

			updated, err := repo.Update(updateData)
			assert.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr == nil {
				assert.NotNil(t, updated)
				assert.Equal(t, updateData, updated)

				var stored ownership
				err = db.First(&stored, updateData.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, updateData.ItemID, stored.ItemID)
				assert.Equal(t, updateData.UserID, stored.UserID)
				assert.Equal(t, updateData.Rentable, stored.Rentable)
				assert.Equal(t, updateData.Memo, stored.Memo)
			}
		})
	}
}

func TestOwnershipRepository_Delete(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(t *testing.T, db *gorm.DB) int
		expectedErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) int {
				createTestItems(t, db, 1)
				model := &ownership{ItemID: 1, UserID: "user1", Rentable: true, Memo: "memo"}
				if err := db.Create(model).Error; err != nil {
					t.Fatalf("failed to create test ownership: %v", err)
				}
				return model.ID
			},
			expectedErr: nil,
		},
		{
			name: "failure: ownership not found",
			setup: func(t *testing.T, db *gorm.DB) int {
				return 9999
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewOwnershipRepository(db)
			id := tc.setup(t, db)

			err := repo.Delete(id)

			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				var model ownership
				err := db.First(&model, id).Error
				assert.Error(t, err)
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			}
		})
	}
}
