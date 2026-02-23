package repository

import (
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
			db := setupTestDB(t)

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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

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
			name: "success: create new item if not exists",
			setup: func(t *testing.T, db *gorm.DB) *domain.Item {
				return &domain.Item{ID: 9999}
			},
			updateItem: &domain.Item{
				Name:        "Non-existent Item",
				Description: "This item does not exist",
				ImgUrl:      "http://example.com/non_existent_image.png",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			repo := NewItemRepository(db)
			existingItem := tc.setup(t, db)
			tc.updateItem.ID = existingItem.ID

			updatedItem, err := repo.Update(tc.updateItem)

			if tc.expectedErr != nil {
				assert.Error(t, err)
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
			db := setupTestDB(t)

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

func TestItemRepository_Search(t *testing.T) {
	type testCase struct {
		name        string
		createItems []*domain.Item
		query       domain.ItemSearchQuery
		expected    []*domain.Item
		wantErr     bool
	}

	testContexts := []struct {
		name      string
		testCases []testCase
	}{

		{
			name: "search by name",
			testCases: []testCase{
				{
					name: "success: empty query",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					wantErr: false,
				},
				{
					name: "success: exact match",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "Test Item 1",
					},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
					},
					wantErr: false,
				},
				{
					name: "success: partial match",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "Test",
					},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
					},
					wantErr: false,
				},
				{
					name: "success: no match",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "Non-existent",
					},
					expected: []*domain.Item{},
					wantErr:  false,
				},
				{
					name: "success: multiple matches",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "Item",
					},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					wantErr: false,
				},
				{
					name: "success: case insensitive match",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "test item 1",
					},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
					},
					wantErr: false,
				},
				{
					name: "success: japanese name",
					createItems: []*domain.Item{
						{Name: "テスト物品1", Description: "なんらかの説明1", ImgUrl: "http://example.com/image1.png"},
						{Name: "テスト物品2", Description: "なんらかの説明2", ImgUrl: "http://example.com/image2.png"},
						{Name: "別の物品", Description: "別の説明", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "テスト物品1",
					},
					expected: []*domain.Item{
						{Name: "テスト物品1", Description: "なんらかの説明1", ImgUrl: "http://example.com/image1.png"},
					},
					wantErr: false,
				},
				{
					name: "success: japanese partial match",
					createItems: []*domain.Item{
						{Name: "テスト物品1", Description: "なんらかの説明1", ImgUrl: "http://example.com/image1.png"},
						{Name: "テスト物品2", Description: "なんらかの説明2", ImgUrl: "http://example.com/image2.png"},
						{Name: "別の物品", Description: "別の説明", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "テスト",
					},
					expected: []*domain.Item{
						{Name: "テスト物品1", Description: "なんらかの説明1", ImgUrl: "http://example.com/image1.png"},
						{Name: "テスト物品2", Description: "なんらかの説明2", ImgUrl: "http://example.com/image2.png"},
					},
					wantErr: false,
				},
			},
		},
		{
			name: "limit, offset",
			testCases: []testCase{
				{
					name: "success: limit results",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Limit: 2,
					},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
					},
				},
				{
					name: "success: limit with offset",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Limit:  1,
						Offset: 1,
					},
					expected: []*domain.Item{
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
					},
					wantErr: false,
				},
				{
					name: "success: offset exceeds total",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Offset: 5,
					},
					expected: []*domain.Item{},
					wantErr:  false,
				},
				{
					name: "success: limit exceeds total",
					createItems: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgUrl: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Limit: 5,
					},
					expected: []*domain.Item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgUrl: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgUrl: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgUrl: "http://example.com/image3.png"},
					},
					wantErr: false,
				},
			},
		},
	}

	for _, tc := range testContexts {
		t.Run(tc.name, func(t *testing.T) {
			for _, c := range tc.testCases {
				t.Run(c.name, func(t *testing.T) {
					db := setupTestDB(t)
					repo := NewItemRepository(db)

					for _, item := range c.createItems {
						_, err := repo.Create(item)
						assert.NoError(t, err)
					}

					results, err := repo.Search(c.query)
					if c.wantErr {
						assert.Error(t, err)
						assert.Nil(t, results)
					} else {
						assert.NoError(t, err)
						assert.Equal(t, len(c.expected), len(results))
						for i := range c.expected {
							assert.Equal(t, c.expected[i].Name, results[i].Name)
							assert.Equal(t, c.expected[i].Description, results[i].Description)
							assert.Equal(t, c.expected[i].ImgUrl, results[i].ImgUrl)
						}
					}
				})
			}
		})
	}
}
