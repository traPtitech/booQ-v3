package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	mock_domain "github.com/traPtitech/booQ-v3/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestItemUseCase_GetItemByID(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(repo *mock_domain.MockItemRepository)
		id           int
		expectedItem *domain.Item
		expectedErr  error
	}{
		{
			name: "success",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Item{
						ID:          1,
						Name:        "Test Item",
						Description: "This is a test item",
						ImgUrl:      "http://example.com/image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "1234567890",
						},
						EquipmentDetail: nil,
					}, nil).
					Times(1)
			},
			id: 1,
			expectedItem: &domain.Item{
				ID:          1,
				Name:        "Test Item",
				Description: "This is a test item",
				ImgUrl:      "http://example.com/image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "1234567890",
				},
				EquipmentDetail: nil,
			},
			expectedErr: nil,
		},
		{
			name: "failure: item not found",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetByID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			id:           2,
			expectedItem: nil,
			expectedErr:  domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemRepo := mock_domain.NewMockItemRepository(ctrl)
			tc.setupMock(mockItemRepo)

			itemUseCase := NewItemUseCase(mockItemRepo)

			item, err := itemUseCase.GetItemByID(tc.id)

			assert.Equal(t, tc.expectedItem, item)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestItemUseCase_GetItemDetailByID(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(repo *mock_domain.MockItemRepository)
		id           int
		expectedItem *domain.ItemDetail
		expectedErr  error
	}{
		{
			name: "success",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetDetailByID(1).
					Return(&domain.ItemDetail{
						Item:  &domain.Item{ID: 1, Name: "Test Item"},
						Tags:  []*domain.Tag{{Name: "tag1"}},
						Likes: []*domain.Like{{ItemID: 1, UserID: "user1"}},
					}, nil).
					Times(1)
			},
			id: 1,
			expectedItem: &domain.ItemDetail{
				Item:  &domain.Item{ID: 1, Name: "Test Item"},
				Tags:  []*domain.Tag{{Name: "tag1"}},
				Likes: []*domain.Like{{ItemID: 1, UserID: "user1"}},
			},
		},
		{
			name: "failure: item not found",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetDetailByID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			id:           2,
			expectedItem: nil,
			expectedErr:  domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemRepo := mock_domain.NewMockItemRepository(ctrl)
			tc.setupMock(mockItemRepo)

			itemUseCase := NewItemUseCase(mockItemRepo)

			item, err := itemUseCase.GetItemDetailByID(tc.id)

			assert.Equal(t, tc.expectedItem, item)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestItemUseCase_CreateItem(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(repo *mock_domain.MockItemRepository)
		inputItem    *domain.Item
		expectedItem *domain.Item
		expectedErr  error
	}{
		{
			name: "success",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					Create(gomock.Any()).
					DoAndReturn(func(item *domain.Item) (*domain.Item, error) {
						item.ID = 1
						return item, nil
					}).
					Times(1)
			},
			inputItem: &domain.Item{
				Name:        "New Item",
				Description: "This is a new item",
				ImgUrl:      "http://example.com/new_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "0987654321",
				},
				EquipmentDetail: nil,
			},
			expectedItem: &domain.Item{
				ID:          1,
				Name:        "New Item",
				Description: "This is a new item",
				ImgUrl:      "http://example.com/new_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "0987654321",
				},
				EquipmentDetail: nil,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemRepo := mock_domain.NewMockItemRepository(ctrl)
			tc.setupMock(mockItemRepo)

			itemUseCase := NewItemUseCase(mockItemRepo)

			createdItem, err := itemUseCase.CreateItem(tc.inputItem)

			assert.Equal(t, tc.expectedItem, createdItem)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestItemUseCase_SearchItems(t *testing.T) {
	testCases := []struct {
		name          string
		setupMock     func(repo *mock_domain.MockItemRepository)
		query         domain.ItemSearchQuery
		expectedItems []*domain.ItemDetail
		expectedErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					Search(domain.ItemSearchQuery{Name: "Test", Limit: 10}).
					Return([]*domain.ItemDetail{
						{
							Item:  &domain.Item{ID: 1, Name: "Test Item"},
							Tags:  []*domain.Tag{{Name: "tag1"}},
							Likes: []*domain.Like{{ItemID: 1, UserID: "user1"}},
						},
					}, nil).
					Times(1)
			},
			query: domain.ItemSearchQuery{
				Name:  "Test",
				Limit: 10,
			},
			expectedItems: []*domain.ItemDetail{
				{
					Item:  &domain.Item{ID: 1, Name: "Test Item"},
					Tags:  []*domain.Tag{{Name: "tag1"}},
					Likes: []*domain.Like{{ItemID: 1, UserID: "user1"}},
				},
			},
		},
		{
			name:      "failure: no limit with offset",
			setupMock: func(repo *mock_domain.MockItemRepository) {},
			query: domain.ItemSearchQuery{
				Offset: 10,
			},
			expectedErr: ErrInvalidSearchQuery,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemRepo := mock_domain.NewMockItemRepository(ctrl)
			tc.setupMock(mockItemRepo)

			itemUseCase := NewItemUseCase(mockItemRepo)

			items, err := itemUseCase.SearchItems(tc.query)

			assert.Equal(t, tc.expectedItems, items)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestItemUseCase_UpdateItem(t *testing.T) {
	testCases := []struct {
		name         string
		setupMock    func(repo *mock_domain.MockItemRepository)
		inputItem    *domain.Item
		expectedItem *domain.Item
		expectedErr  error
	}{
		{
			name: "success",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Item{
						ID:          1,
						Name:        "Test Item",
						Description: "This is a test item",
						ImgUrl:      "http://example.com/image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "1234567890",
						},
						EquipmentDetail: nil,
					}, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any()).
					Return(&domain.Item{
						ID:          1,
						Name:        "Updated Item",
						Description: "This is an updated item",
						ImgUrl:      "http://example.com/updated_image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "1234567890",
						},
						EquipmentDetail: nil,
					}, nil).
					Times(1)
			},
			inputItem: &domain.Item{
				ID:          1,
				Name:        "Updated Item",
				Description: "This is an updated item",
				ImgUrl:      "http://example.com/updated_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "1234567890",
				},
				EquipmentDetail: nil,
			},
			expectedItem: &domain.Item{
				ID:          1,
				Name:        "Updated Item",
				Description: "This is an updated item",
				ImgUrl:      "http://example.com/updated_image.png",
				BookDetail: &domain.BookDetail{
					ISBNCode: "1234567890",
				},
				EquipmentDetail: nil,
			},
			expectedErr: nil,
		},
		{
			name: "failure: item not found",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetByID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			inputItem: &domain.Item{
				ID:              2,
				Name:            "Nonexistent Item",
				Description:     "This item does not exist",
				ImgUrl:          "http://example.com/nonexistent_image.png",
				BookDetail:      nil,
				EquipmentDetail: nil,
			},
			expectedItem: nil,
			expectedErr:  domain.ErrNotFound,
		},
		{
			name: "failure: cannot change book detail",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Item{
						ID:          1,
						Name:        "Test Item",
						Description: "This is a test item",
						ImgUrl:      "http://example.com/image.png",
						BookDetail: &domain.BookDetail{
							ISBNCode: "1234567890",
						},
						EquipmentDetail: nil,
					}, nil).
					Times(1)
			},
			inputItem: &domain.Item{
				ID:              1,
				Name:            "Updated Item",
				Description:     "This is an updated item",
				ImgUrl:          "http://example.com/updated_image.png",
				BookDetail:      nil, // book detail removed
				EquipmentDetail: nil,
			},
			expectedItem: nil,
			expectedErr:  ErrUpdateNotAllowed,
		},
		{
			name: "failure: cannot change equipment detail",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Item{
						ID:          1,
						Name:        "Test Item",
						Description: "This is a test item",
						ImgUrl:      "http://example.com/image.png",
						BookDetail:  nil,
						EquipmentDetail: &domain.EquipmentDetail{
							Count:    5,
							CountMax: 10,
						},
					}, nil).
					Times(1)
			},
			inputItem: &domain.Item{
				ID:              1,
				Name:            "Updated Item",
				Description:     "This is an updated item",
				ImgUrl:          "http://example.com/updated_image.png",
				BookDetail:      nil,
				EquipmentDetail: nil, // equipment detail removed
			},
			expectedItem: nil,
			expectedErr:  ErrUpdateNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemRepo := mock_domain.NewMockItemRepository(ctrl)
			tc.setupMock(mockItemRepo)

			itemUseCase := NewItemUseCase(mockItemRepo)

			updatedItem, err := itemUseCase.UpdateItem(tc.inputItem)

			assert.Equal(t, tc.expectedItem, updatedItem)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestItemUseCase_DeleteItem(t *testing.T) {
	testCases := []struct {
		name        string
		setupMock   func(repo *mock_domain.MockItemRepository)
		id          int
		expectedErr error
	}{
		{
			name: "success",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					Delete(1).
					Return(nil).
					Times(1)
			},
			id:          1,
			expectedErr: nil,
		},
		{
			name: "failure: item not found",
			setupMock: func(repo *mock_domain.MockItemRepository) {
				repo.EXPECT().
					Delete(2).
					Return(domain.ErrNotFound).
					Times(1)
			},
			id:          2,
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockItemRepo := mock_domain.NewMockItemRepository(ctrl)
			tc.setupMock(mockItemRepo)

			itemUseCase := NewItemUseCase(mockItemRepo)

			err := itemUseCase.DeleteItem(tc.id)

			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}
