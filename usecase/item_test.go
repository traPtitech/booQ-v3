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
					Return(nil, domain.ErrItemNotFound).
					Times(1)
			},
			id:           2,
			expectedItem: nil,
			expectedErr:  domain.ErrItemNotFound,
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
