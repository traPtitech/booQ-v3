package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

func TestItemRepository_GetByID(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func(t *testing.T, db *gorm.DB) int
		verify func(t *testing.T, item *domain.Item, err error)
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) int {
				item := &item{
					Name:        "Test Item",
					Description: "This is a test item",
					ImgURL:      "http://example.com/image.png",
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item: %v", err)
				}
				return item.ID
			},
			verify: func(t *testing.T, item *domain.Item, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, "Test Item", item.Name)
				assert.Equal(t, "This is a test item", item.Description)
				assert.Equal(t, "http://example.com/image.png", item.ImgUrl)
			},
		},
		{
			name: "failure: item not found",
			setup: func(t *testing.T, db *gorm.DB) int {
				return 9999 // not existing ID
			},
			verify: func(t *testing.T, item *domain.Item, err error) {
				assert.Error(t, err)
				assert.ErrorIs(t, err, domain.ErrItemNotFound)
				assert.Nil(t, item)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := setupTestDB(ctx, t)

			repo := NewItemRepository(db)
			id := tc.setup(t, db)

			item, err := repo.GetByID(id)
			tc.verify(t, item, err)
		})
	}
}

func TestItemRepository_Create(t *testing.T) {
	testCases := []struct {
		name    string
		item    *domain.Item
		wantErr bool
	}{
		{
			name: "success",
			item: &domain.Item{
				Name:        "New Item",
				Description: "This is a new item",
				ImgUrl:      "http://example.com/new_image.png",
			},
			wantErr: false,
		},
		{
			name: "failure: missing name",
			item: &domain.Item{
				Description: "This item has no name",
				ImgUrl:      "http://example.com/no_name_image.png",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := setupTestDB(ctx, t)

			repo := NewItemRepository(db)
			createdItem, err := repo.Create(tc.item)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, createdItem)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, createdItem)
				assert.Equal(t, tc.item.Name, createdItem.Name)
				assert.Equal(t, tc.item.Description, createdItem.Description)
				assert.Equal(t, tc.item.ImgUrl, createdItem.ImgUrl)
			}
		})
	}
}

func TestItemRepository_Update(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(t *testing.T, db *gorm.DB) *domain.Item
		updateItem  *domain.Item
		expectedErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) *domain.Item {
				item := &item{
					Name:        "Original Item",
					Description: "This is the original item",
					ImgURL:      "http://example.com/original_image.png",
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item: %v", err)
				}
				return item.toDomain()
			},
			updateItem: &domain.Item{
				Name:        "Updated Item",
				Description: "This is the updated item",
				ImgUrl:      "http://example.com/updated_image.png",
			},
			expectedErr: nil,
		},
		{
			name: "failure: item not found",
			setup: func(t *testing.T, db *gorm.DB) *domain.Item {
				return &domain.Item{ID: 9999}
			},
			updateItem: &domain.Item{
				Name:        "Non-existent Item",
				Description: "This item does not exist",
				ImgUrl:      "http://example.com/non_existent_image.png",
			},
			expectedErr: domain.ErrItemNotFound,
		},
		{
			name: "failure: missing name",
			setup: func(t *testing.T, db *gorm.DB) *domain.Item {
				item := &item{
					Name:        "Item to Update",
					Description: "This item will be updated",
					ImgURL:      "http://example.com/item_to_update_image.png",
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item: %v", err)
				}
				return item.toDomain()
			},
			updateItem: &domain.Item{
				Description: "This item has no name",
				ImgUrl:      "http://example.com/no_name_update_image.png",
			},
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := setupTestDB(ctx, t)

			repo := NewItemRepository(db)
			existingItem := tc.setup(t, db)
			tc.updateItem.ID = existingItem.ID

			updatedItem, err := repo.Update(tc.updateItem)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				if !errors.Is(tc.expectedErr, assert.AnError) {
					assert.ErrorIs(t, err, tc.expectedErr)
				}
				assert.Nil(t, updatedItem)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedItem)
				assert.Equal(t, tc.updateItem.Name, updatedItem.Name)
				assert.Equal(t, tc.updateItem.Description, updatedItem.Description)
				assert.Equal(t, tc.updateItem.ImgUrl, updatedItem.ImgUrl)
			}
		})
	}
}

func TestItemRepository_Delete(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(t *testing.T, db *gorm.DB) int
		expectedErr error
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) int {
				item := &item{
					Name:        "Item to Delete",
					Description: "This item will be deleted",
					ImgURL:      "http://example.com/item_to_delete_image.png",
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item: %v", err)
				}
				return item.ID
			},
			expectedErr: nil,
		},
		{
			name: "failure: item not found",
			setup: func(t *testing.T, db *gorm.DB) int {
				return 9999
			},
			expectedErr: domain.ErrItemNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := setupTestDB(ctx, t)

			repo := NewItemRepository(db)
			id := tc.setup(t, db)

			err := repo.Delete(id)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)

				var item item
				err = db.First(&item, id).Error
				assert.Error(t, err)
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			}
		})
	}
}
