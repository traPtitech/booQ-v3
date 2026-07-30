package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	mock_domain "github.com/traPtitech/booQ-v3/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestLikeUseCase_AddLike(t *testing.T) {
	testCases := []struct {
		name        string
		itemID      int
		userID      string
		setupMock   func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository)
		expectedErr error
	}{
		{
			name:   "success",
			itemID: 1,
			userID: "user1",
			setupMock: func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository) {
				itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
				likeRepo.EXPECT().Exists(1, "user1").Return(false, nil)
				likeRepo.EXPECT().Create(&domain.Like{ItemID: 1, UserID: "user1"}).Return(nil)
			},
		},
		{
			name:   "failure: item not found",
			itemID: 999,
			userID: "user1",
			setupMock: func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository) {
				itemRepo.EXPECT().GetByID(999).Return(nil, domain.ErrNotFound)
			},
			expectedErr: domain.ErrNotFound,
		},
		{
			name:   "failure: already liked",
			itemID: 1,
			userID: "user1",
			setupMock: func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository) {
				itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
				likeRepo.EXPECT().Exists(1, "user1").Return(true, nil)
			},
			expectedErr: ErrAlreadyLiked,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			itemRepo := mock_domain.NewMockItemRepository(ctrl)
			likeRepo := mock_domain.NewMockLikeRepository(ctrl)
			tc.setupMock(itemRepo, likeRepo)

			u := NewLikeUseCase(likeRepo, itemRepo)
			err := u.AddLike(tc.itemID, tc.userID)

			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestLikeUseCase_RemoveLike(t *testing.T) {
	testCases := []struct {
		name        string
		itemID      int
		userID      string
		setupMock   func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository)
		expectedErr error
	}{
		{
			name:   "success",
			itemID: 1,
			userID: "user1",
			setupMock: func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository) {
				itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
				likeRepo.EXPECT().Exists(1, "user1").Return(true, nil)
				likeRepo.EXPECT().Delete(1, "user1").Return(nil)
			},
		},
		{
			name:   "failure: item not found",
			itemID: 999,
			userID: "user1",
			setupMock: func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository) {
				itemRepo.EXPECT().GetByID(999).Return(nil, domain.ErrNotFound)
			},
			expectedErr: domain.ErrNotFound,
		},
		{
			name:   "failure: not liked",
			itemID: 1,
			userID: "user1",
			setupMock: func(itemRepo *mock_domain.MockItemRepository, likeRepo *mock_domain.MockLikeRepository) {
				itemRepo.EXPECT().GetByID(1).Return(&domain.Item{ID: 1}, nil)
				likeRepo.EXPECT().Exists(1, "user1").Return(false, nil)
			},
			expectedErr: ErrNotLiked,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			itemRepo := mock_domain.NewMockItemRepository(ctrl)
			likeRepo := mock_domain.NewMockLikeRepository(ctrl)
			tc.setupMock(itemRepo, likeRepo)

			u := NewLikeUseCase(likeRepo, itemRepo)
			err := u.RemoveLike(tc.itemID, tc.userID)

			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
