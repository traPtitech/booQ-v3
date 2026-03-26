package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	mock_domain "github.com/traPtitech/booQ-v3/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestTagUseCase_ReplaceByItemID(t *testing.T) {
	testCases := []struct {
		name        string
		itemID      int
		tags        []string
		setupMock   func(itemRepo *mock_domain.MockItemRepository, tagRepo *mock_domain.MockTagRepository)
		expectedErr error
	}{
		{
			name:   "success",
			itemID: 1,
			tags:   []string{"go", "book"},
			setupMock: func(itemRepo *mock_domain.MockItemRepository, tagRepo *mock_domain.MockTagRepository) {
				itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
				tagRepo.EXPECT().ReplaceByItemID(1, []string{"go", "book"}).Return(nil)
			},
		},
		{
			name:   "failure: item not found",
			itemID: 999,
			tags:   []string{"go"},
			setupMock: func(itemRepo *mock_domain.MockItemRepository, tagRepo *mock_domain.MockTagRepository) {
				itemRepo.EXPECT().GetByID(999).Return(nil, domain.ErrNotFound)
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			itemRepo := mock_domain.NewMockItemRepository(ctrl)
			tagRepo := mock_domain.NewMockTagRepository(ctrl)
			tc.setupMock(itemRepo, tagRepo)

			u := NewTagUseCase(tagRepo, itemRepo)
			err := u.ReplaceByItemID(tc.itemID, tc.tags)

			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
