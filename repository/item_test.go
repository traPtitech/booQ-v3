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
				assert.ErrorIs(t, err, domain.ErrNotFound)
				assert.Nil(t, item)
			},
		},
		{
			name: "success: item with book detail",
			setup: func(t *testing.T, db *gorm.DB) int {
				item := &item{
					Name:        "Book Item",
					Description: "This is a book item",
					ImgURL:      "http://example.com/book_image.png",
					Book: &book{
						ISBNCode: "1234567890123",
					},
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item with book detail: %v", err)
				}
				return item.ID
			},
			verify: func(t *testing.T, item *domain.Item, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, "Book Item", item.Name)
				assert.Equal(t, "This is a book item", item.Description)
				assert.Equal(t, "http://example.com/book_image.png", item.ImgUrl)
				assert.NotNil(t, item.BookDetail)
				assert.Equal(t, "1234567890123", item.BookDetail.ISBNCode)
				assert.Nil(t, item.EquipmentDetail)
			},
		},
		{
			name: "success: item with equipment detail",
			setup: func(t *testing.T, db *gorm.DB) int {
				item := &item{
					Name:        "Equipment Item",
					Description: "This is an equipment item",
					ImgURL:      "http://example.com/equipment_image.png",
					Equipment: &equipment{
						Count:    5,
						CountMax: 10,
					},
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item with equipment detail: %v", err)
				}
				return item.ID
			},
			verify: func(t *testing.T, item *domain.Item, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, "Equipment Item", item.Name)
				assert.Equal(t, "This is an equipment item", item.Description)
				assert.Equal(t, "http://example.com/equipment_image.png", item.ImgUrl)
				assert.NotNil(t, item.EquipmentDetail)
				assert.Equal(t, 5, item.EquipmentDetail.Count)
				assert.Equal(t, 10, item.EquipmentDetail.CountMax)
				assert.Nil(t, item.BookDetail)
			},
		},
		{
			name: "success: item with both book and equipment detail",
			setup: func(t *testing.T, db *gorm.DB) int {
				item := &item{
					Name:        "Complex Item",
					Description: "This item has both book and equipment details",
					ImgURL:      "http://example.com/complex_image.png",
					Book: &book{
						ISBNCode: "9876543210123",
					},
					Equipment: &equipment{
						Count:    3,
						CountMax: 5,
					},
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item with both details: %v", err)
				}
				return item.ID
			},
			verify: func(t *testing.T, item *domain.Item, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.Equal(t, "Complex Item", item.Name)
				assert.Equal(t, "This item has both book and equipment details", item.Description)
				assert.Equal(t, "http://example.com/complex_image.png", item.ImgUrl)
				assert.NotNil(t, item.BookDetail)
				assert.Equal(t, "9876543210123", item.BookDetail.ISBNCode)
				assert.NotNil(t, item.EquipmentDetail)
				assert.Equal(t, 3, item.EquipmentDetail.Count)
				assert.Equal(t, 5, item.EquipmentDetail.CountMax)
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

func TestItemModel_PreloadLikesAndTags(t *testing.T) {
	db := setupTestDB(t)

	model := &item{
		Name:        "Item with likes and tags",
		Description: "preload test",
		ImgURL:      "http://example.com/preload.png",
		Likes: []like{
			{UserID: "user-a"},
			{UserID: "user-b"},
		},
		Tags: []tag{
			{Name: "go"},
			{Name: "book"},
		},
	}
	if err := db.Create(model).Error; err != nil {
		t.Fatalf("Failed to create test item with likes and tags: %v", err)
	}

	var got item
	err := db.Preload("Likes").Preload("Tags").First(&got, model.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, model.ID, got.ID)

	assert.Len(t, got.Likes, 2)
	assert.ElementsMatch(t, []string{"user-a", "user-b"}, []string{got.Likes[0].UserID, got.Likes[1].UserID})
	for _, like := range got.Likes {
		assert.Equal(t, got.ID, like.ItemID)
	}

	assert.Len(t, got.Tags, 2)
	assert.ElementsMatch(t, []string{"go", "book"}, []string{got.Tags[0].Name, got.Tags[1].Name})
	for _, tag := range got.Tags {
		assert.Equal(t, got.ID, tag.ItemID)
	}
}

func TestItemRepository_GetDetailByID(t *testing.T) {
	testCases := []struct {
		name   string
		setup  func(t *testing.T, db *gorm.DB) int
		verify func(t *testing.T, item *domain.ItemDetail, err error)
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) int {
				model := &item{
					Name:        "Detailed Item",
					Description: "detail test",
					ImgURL:      "http://example.com/detail.png",
					Tags: []tag{
						{Name: "go"},
						{Name: "book"},
					},
					Likes: []like{
						{UserID: "user-a"},
						{UserID: "user-b"},
					},
					Ownership: []ownership{
						{
							UserID:   "owner1",
							Rentable: true,
							Transaction: []transaction{
								{UserID: "borrower1", Status: "requested", Purpose: "read"},
							},
						},
						{
							UserID:   "owner2",
							Rentable: false,
							Transaction: []transaction{
								{UserID: "borrower2", Status: "borrowed", Purpose: "research"},
							},
						},
					},
				}
				if err := db.Create(model).Error; err != nil {
					t.Fatalf("Failed to create detailed item: %v", err)
				}
				return model.ID
			},
			verify: func(t *testing.T, item *domain.ItemDetail, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, item)
				assert.NotNil(t, item.Item)
				assert.Equal(t, "Detailed Item", item.Item.Name)
				assert.Len(t, item.Tags, 2)
				assert.Len(t, item.Likes, 2)
				assert.Len(t, item.Ownerships, 2)
				assert.NotNil(t, item.Ownerships[0].Ownership)
				assert.Equal(t, "owner1", item.Ownerships[0].Ownership.UserID)
				assert.Len(t, item.Ownerships[0].Transactions, 1)
				assert.Equal(t, "borrower1", item.Ownerships[0].Transactions[0].UserID)
				assert.NotNil(t, item.Ownerships[1].Ownership)
				assert.Equal(t, "owner2", item.Ownerships[1].Ownership.UserID)
				assert.Len(t, item.Ownerships[1].Transactions, 1)
				assert.Equal(t, "borrower2", item.Ownerships[1].Transactions[0].UserID)
			},
		},
		{
			name:  "failure: item not found",
			setup: func(t *testing.T, db *gorm.DB) int { return 9999 },
			verify: func(t *testing.T, item *domain.ItemDetail, err error) {
				assert.ErrorIs(t, err, domain.ErrNotFound)
				assert.Nil(t, item)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			repo := NewItemRepository(db)
			id := tc.setup(t, db)

			item, err := repo.GetDetailByID(id)
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
			name: "success: item with book detail",
			item: &domain.Item{
				Name:        "Book Item",
				Description: "This is a book item",
				ImgUrl:      "http://example.com/book_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "1234567890123",
				},
			},
			wantErr: false,
		},
		{
			name: "success: item with equipment detail",
			item: &domain.Item{
				Name:        "Equipment Item",
				Description: "This is an equipment item",
				ImgUrl:      "http://example.com/equipment_image.png",
				EquipmentDetail: &domain.EquipmentDetail{
					Count:    5,
					CountMax: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "success: item with both book and equipment detail",
			item: &domain.Item{
				Name:        "Complex Item",
				Description: "This item has both book and equipment details",
				ImgUrl:      "http://example.com/complex_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "9876543210123",
				},
				EquipmentDetail: &domain.EquipmentDetail{
					Count:    3,
					CountMax: 5,
				},
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
				assert.NotEqual(t, 0, createdItem.ID)

				if tc.item.BookDetail != nil {
					assert.NotNil(t, createdItem.BookDetail)
					assert.Equal(t, tc.item.BookDetail.ISBNCode, createdItem.BookDetail.ISBNCode)
				} else {
					assert.Nil(t, createdItem.BookDetail)
				}

				if tc.item.EquipmentDetail != nil {
					assert.NotNil(t, createdItem.EquipmentDetail)
					assert.Equal(t, tc.item.EquipmentDetail.Count, createdItem.EquipmentDetail.Count)
					assert.Equal(t, tc.item.EquipmentDetail.CountMax, createdItem.EquipmentDetail.CountMax)
				} else {
					assert.Nil(t, createdItem.EquipmentDetail)
				}

				var item item
				err = db.First(&item, createdItem.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tc.item.Name, item.Name)
				assert.Equal(t, tc.item.Description, item.Description)
				assert.Equal(t, tc.item.ImgUrl, item.ImgURL)

				if tc.item.BookDetail != nil {
					var book book
					err = db.First(&book, "item_id = ?", createdItem.ID).Error
					assert.NoError(t, err)
					assert.Equal(t, tc.item.BookDetail.ISBNCode, book.ISBNCode)
				}

				if tc.item.EquipmentDetail != nil {
					var equipment equipment
					err = db.First(&equipment, "item_id = ?", createdItem.ID).Error
					assert.NoError(t, err)
					assert.Equal(t, tc.item.EquipmentDetail.Count, equipment.Count)
					assert.Equal(t, tc.item.EquipmentDetail.CountMax, equipment.CountMax)
				}
			}
		})
	}
}

func TestItemRepository_CreateBatch(t *testing.T) {
	testCases := []struct {
		name    string
		items   []*domain.Item
		wantErr bool
	}{
		{
			name: "success: batch create multiple items",
			items: []*domain.Item{
				{Name: "Batch Item 1", Description: "Description 1", ImgUrl: "http://example.com/1.png"},
				{Name: "Batch Item 2", Description: "Description 2", ImgUrl: "http://example.com/2.png"},
				{Name: "Batch Item 3", Description: "Description 3", ImgUrl: "http://example.com/3.png"},
			},
			wantErr: false,
		},
		{
			name:    "success: batch create empty slice",
			items:   []*domain.Item{},
			wantErr: false,
		},
		{
			name: "success: batch create items with details",
			items: []*domain.Item{
				{
					Name:        "Batch Book Item",
					Description: "This is a batch book item",
					ImgUrl:      "http://example.com/book.png",
					BookDetail: &domain.BookDetail{
						ISBNCode: "1234567890123",
					},
				},
				{
					Name:        "Batch Equipment Item",
					Description: "This is a batch equipment item",
					ImgUrl:      "http://example.com/equipment.png",
					EquipmentDetail: &domain.EquipmentDetail{
						Count:    5,
						CountMax: 10,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)

			repo := NewItemRepository(db)
			createdItems, err := repo.CreateBatch(tc.items)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, createdItems)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.items), len(createdItems))
				for i, expected := range tc.items {
					assert.Equal(t, expected.Name, createdItems[i].Name)
					assert.Equal(t, expected.Description, createdItems[i].Description)
					assert.Equal(t, expected.ImgUrl, createdItems[i].ImgUrl)
					assert.NotEqual(t, 0, createdItems[i].ID)

					if expected.BookDetail != nil {
						assert.NotNil(t, createdItems[i].BookDetail)
						assert.Equal(t, expected.BookDetail.ISBNCode, createdItems[i].BookDetail.ISBNCode)
					} else {
						assert.Nil(t, createdItems[i].BookDetail)
					}

					if expected.EquipmentDetail != nil {
						assert.NotNil(t, createdItems[i].EquipmentDetail)
						assert.Equal(t, expected.EquipmentDetail.Count, createdItems[i].EquipmentDetail.Count)
						assert.Equal(t, expected.EquipmentDetail.CountMax, createdItems[i].EquipmentDetail.CountMax)
					} else {
						assert.Nil(t, createdItems[i].EquipmentDetail)
					}

					var item item
					err = db.First(&item, createdItems[i].ID).Error
					assert.NoError(t, err)
					assert.Equal(t, expected.Name, item.Name)
				}
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
		{
			name: "success: update item with book detail",
			setup: func(t *testing.T, db *gorm.DB) *domain.Item {
				item := &item{
					Name:        "Book Item",
					Description: "This is a book item",
					ImgURL:      "http://example.com/book_image.png",
					Book: &book{
						ISBNCode: "1234567890123",
					},
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item with book detail: %v", err)
				}
				return item.toDomain()
			},
			updateItem: &domain.Item{
				Name:        "Updated Book Item",
				Description: "This is an updated book item",
				ImgUrl:      "http://example.com/updated_book_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "9876543210123",
				},
			},
			expectedErr: nil,
		},
		{
			name: "success: update item with equipment detail",
			setup: func(t *testing.T, db *gorm.DB) *domain.Item {
				item := &item{
					Name:        "Equipment Item",
					Description: "This is an equipment item",
					ImgURL:      "http://example.com/equipment_image.png",
					Equipment: &equipment{
						Count:    5,
						CountMax: 10,
					},
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item with equipment detail: %v", err)
				}
				return item.toDomain()
			},
			updateItem: &domain.Item{
				Name:        "Updated Equipment Item",
				Description: "This is an updated equipment item",
				ImgUrl:      "http://example.com/updated_equipment_image.png",
				EquipmentDetail: &domain.EquipmentDetail{
					Count:    3,
					CountMax: 5,
				},
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
				assert.Equal(t, existingItem.ID, updatedItem.ID)

				if tc.updateItem.BookDetail != nil {
					assert.NotNil(t, updatedItem.BookDetail)
					assert.Equal(t, tc.updateItem.BookDetail.ISBNCode, updatedItem.BookDetail.ISBNCode)
				} else {
					assert.Nil(t, updatedItem.BookDetail)
				}

				if tc.updateItem.EquipmentDetail != nil {
					assert.NotNil(t, updatedItem.EquipmentDetail)
					assert.Equal(t, tc.updateItem.EquipmentDetail.Count, updatedItem.EquipmentDetail.Count)
					assert.Equal(t, tc.updateItem.EquipmentDetail.CountMax, updatedItem.EquipmentDetail.CountMax)
				} else {
					assert.Nil(t, updatedItem.EquipmentDetail)
				}

				var item item
				err = db.First(&item, existingItem.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tc.updateItem.Name, item.Name)
				assert.Equal(t, tc.updateItem.Description, item.Description)
				assert.Equal(t, tc.updateItem.ImgUrl, item.ImgURL)
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
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "success: delete item detail on cascade",
			setup: func(t *testing.T, db *gorm.DB) int {
				item := &item{
					Name:        "Item with Detail to Delete",
					Description: "This item has details and will be deleted",
					ImgURL:      "http://example.com/item_with_detail_to_delete_image.png",
					Book: &book{
						ISBNCode: "1234567890123",
					},
					Equipment: &equipment{
						Count:    5,
						CountMax: 10,
					},
				}
				if err := db.Create(item).Error; err != nil {
					t.Fatalf("Failed to create test item with details: %v", err)
				}
				return item.ID
			},
			expectedErr: nil,
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

				var book book
				err = db.First(&book, "item_id = ?", id).Error
				assert.Error(t, err)
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

				var equipment equipment
				err = db.First(&equipment, "item_id = ?", id).Error
				assert.Error(t, err)
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			}
		})
	}
}

// TODO: expectedもItemDetailにすべき
func TestItemRepository_Search(t *testing.T) {
	type testCase struct {
		name        string
		createItems []*item
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgURL: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Name: "Non-existent",
					},
					expected: []*domain.Item{},
					wantErr:  false,
				},
				{
					name: "success: multiple matches",
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Another Item", Description: "This is another item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "テスト物品1", Description: "なんらかの説明1", ImgURL: "http://example.com/image1.png"},
						{Name: "テスト物品2", Description: "なんらかの説明2", ImgURL: "http://example.com/image2.png"},
						{Name: "別の物品", Description: "別の説明", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "テスト物品1", Description: "なんらかの説明1", ImgURL: "http://example.com/image1.png"},
						{Name: "テスト物品2", Description: "なんらかの説明2", ImgURL: "http://example.com/image2.png"},
						{Name: "別の物品", Description: "別の説明", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgURL: "http://example.com/image3.png"},
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
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgURL: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Limit:  2,
						Offset: 5,
					},
					expected: []*domain.Item{},
					wantErr:  false,
				},
				{
					name: "success: limit exceeds total",
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgURL: "http://example.com/image3.png"},
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
				{
					name: "failure: no limit with offset",
					createItems: []*item{
						{Name: "Test Item 1", Description: "This is the first test item", ImgURL: "http://example.com/image1.png"},
						{Name: "Test Item 2", Description: "This is the second test item", ImgURL: "http://example.com/image2.png"},
						{Name: "Test Item 3", Description: "This is the third test item", ImgURL: "http://example.com/image3.png"},
					},
					query: domain.ItemSearchQuery{
						Offset: 1,
					},
					expected: nil,
					wantErr:  true,
				},
			},
		},
		{
			name: "search by user id",
			testCases: []testCase{
				{
					name: "success: search by owner user id",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Ownership: []ownership{
								{UserID: "owner1"},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Ownership: []ownership{
								{UserID: "owner2"},
							},
						},
						{
							Name:        "Item 3",
							Description: "This is item 3",
							ImgURL:      "http://example.com/image3.png",
							Ownership: []ownership{
								{UserID: "owner1"},
							},
						},
					},
					query: domain.ItemSearchQuery{
						UserID: "owner1",
					},
					expected: []*domain.Item{
						{Name: "Item 1", Description: "This is item 1", ImgUrl: "http://example.com/image1.png"},
						{Name: "Item 3", Description: "This is item 3", ImgUrl: "http://example.com/image3.png"},
					},
					wantErr: false,
				},
			},
		},
		{
			name: "search by borrower id",
			testCases: []testCase{
				{
					name: "success: search by borrower user id",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Ownership: []ownership{
								{
									UserID:   "owner1",
									Rentable: true,
									Transaction: []transaction{
										{UserID: "borrower1", Status: domain.BorrowingStatusBorrowed.ToString(), Purpose: "read"},
									},
								},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Ownership: []ownership{
								{
									UserID:   "owner2",
									Rentable: true,
									Transaction: []transaction{
										{UserID: "borrower2", Status: domain.BorrowingStatusBorrowed.ToString(), Purpose: "research"},
									},
								},
							},
						},
						{
							Name:        "Item 3",
							Description: "This is item 3",
							ImgURL:      "http://example.com/image3.png",
							Ownership: []ownership{
								{
									UserID:   "owner3",
									Rentable: true,
									Transaction: []transaction{
										{UserID: "borrower1", Status: domain.BorrowingStatusBorrowed.ToString(), Purpose: "read"},
									},
								},
							},
						},
					},
					query: domain.ItemSearchQuery{
						UserID: "borrower1",
					},
					expected: []*domain.Item{
						{Name: "Item 1", Description: "This is item 1", ImgUrl: "http://example.com/image1.png"},
						{Name: "Item 3", Description: "This is item 3", ImgUrl: "http://example.com/image3.png"},
					},
					wantErr: false,
				},
				{
					name: "success: get only status is borrowed",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Ownership: []ownership{
								{
									UserID:   "owner1",
									Rentable: true,
									Transaction: []transaction{
										{UserID: "borrower1", Status: domain.BorrowingStatusBorrowed.ToString(), Purpose: "read"},
										{UserID: "borrower2", Status: domain.BorrowingStatusReturned.ToString(), Purpose: "research"},
									},
								},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Ownership: []ownership{
								{
									UserID:   "owner2",
									Rentable: true,
									Transaction: []transaction{
										{UserID: "borrower1", Status: domain.BorrowingStatusRequested.ToString(), Purpose: "read"},
										{UserID: "borrower2", Status: domain.BorrowingStatusBorrowed.ToString(), Purpose: "research"},
									},
								},
							},
						},
					},
					query: domain.ItemSearchQuery{
						UserID: "borrower1",
					},
					expected: []*domain.Item{
						{Name: "Item 1", Description: "This is item 1", ImgUrl: "http://example.com/image1.png"},
					},
					wantErr: false,
				},
			},
		},
		{
			name: "search by tag",
			testCases: []testCase{
				{
					name: "success: search by tag",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Tags: []tag{
								{Name: "Tag1"},
								{Name: "Tag2"},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Tags: []tag{
								{Name: "Tag2"},
								{Name: "Tag3"},
							},
						},
						{
							Name:        "Item 3",
							Description: "This is item 3",
							ImgURL:      "http://example.com/image3.png",
							Tags: []tag{
								{Name: "Tag3"},
								{Name: "Tag4"},
							},
						},
					},
					query: domain.ItemSearchQuery{
						Tag: []string{
							"Tag2",
						},
					},
					expected: []*domain.Item{
						{Name: "Item 1", Description: "This is item 1", ImgUrl: "http://example.com/image1.png"},
						{Name: "Item 2", Description: "This is item 2", ImgUrl: "http://example.com/image2.png"},
					},
					wantErr: false,
				},
				{
					name: "success: search by non-existent tag",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Tags: []tag{
								{Name: "Tag1"},
								{Name: "Tag2"},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Tags: []tag{
								{Name: "Tag2"},
								{Name: "Tag3"},
							},
						},
					},
					query: domain.ItemSearchQuery{
						Tag: []string{
							"NonExistentTag",
						},
					},
					expected: []*domain.Item{},
					wantErr:  false,
				},
				{
					name: "success: search by multiple tags",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Tags: []tag{
								{Name: "Tag1"},
								{Name: "Tag2"},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Tags: []tag{
								{Name: "Tag2"},
								{Name: "Tag3"},
							},
						},
					},
					query: domain.ItemSearchQuery{
						Tag: []string{
							"Tag1",
							"Tag2",
						},
					},
					expected: []*domain.Item{
						{Name: "Item 1", Description: "This is item 1", ImgUrl: "http://example.com/image1.png"},
					},
					wantErr: false,
				},
				{
					name: "success: search with tagExclude",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Tags: []tag{
								{Name: "Tag1"},
								{Name: "Tag2"},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Tags: []tag{
								{Name: "Tag2"},
								{Name: "Tag3"},
							},
						},
					},
					query: domain.ItemSearchQuery{
						TagExclude: []string{
							"Tag2",
						},
					},
					expected: []*domain.Item{
						{Name: "Item 2", Description: "This is item 2", ImgUrl: "http://example.com/image2.png"},
					},
					wantErr: false,
				},
				{
					name: "success: search with tag and tagExclude",
					createItems: []*item{
						{
							Name:        "Item 1",
							Description: "This is item 1",
							ImgURL:      "http://example.com/image1.png",
							Tags: []tag{
								{Name: "Tag1"},
								{Name: "Tag2"},
							},
						},
						{
							Name:        "Item 2",
							Description: "This is item 2",
							ImgURL:      "http://example.com/image2.png",
							Tags: []tag{
								{Name: "Tag2"},
								{Name: "Tag3"},
							},
						},
						{
							Name:        "Item 3",
							Description: "This is item 3",
							ImgURL:      "http://example.com/image3.png",
							Tags: []tag{
								{Name: "Tag1"},
								{Name: "Tag3"},
							},
						},
					},
					query: domain.ItemSearchQuery{
						Tag: []string{
							"Tag1",
						},
						TagExclude: []string{
							"Tag2",
						},
					},
					expected: []*domain.Item{
						{Name: "Item 3", Description: "This is item 3", ImgUrl: "http://example.com/image3.png"},
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
						err := db.Create(item).Error
						assert.NoError(t, err)
					}

					results, err := repo.Search(c.query)
					if c.wantErr {
						assert.Error(t, err)
						assert.Nil(t, results)
					} else {
						if assert.NoError(t, err) {
							assert.Equal(t, len(c.expected), len(results))
							for i := range c.expected {
								assert.NotNil(t, results[i].Item)
								assert.Equal(t, c.expected[i].Name, results[i].Item.Name)
								assert.Equal(t, c.expected[i].Description, results[i].Item.Description)
								assert.Equal(t, c.expected[i].ImgUrl, results[i].Item.ImgUrl)
							}
						}
					}
				})
			}
		})
	}
}
