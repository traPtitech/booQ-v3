package repository

import (
	"context"
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
