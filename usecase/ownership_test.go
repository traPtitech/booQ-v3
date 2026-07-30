package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	mock_domain "github.com/traPtitech/booQ-v3/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestOwnershipUseCase_GetByItemID(t *testing.T) {
	testCases := []struct {
		name               string
		itemID             int
		setupMock          func(repo *mock_domain.MockOwnershipRepository)
		expectedOwnerships []*domain.Ownership
		expectedErr        error
	}{
		{
			name:   "success",
			itemID: 1,
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByItemID(1).
					Return([]*domain.Ownership{
						{ID: 1, ItemID: 1, UserID: "user1", Rentable: true, Memo: "memo1"},
						{ID: 2, ItemID: 1, UserID: "user2", Rentable: false, Memo: "memo2"},
					}, nil).
					Times(1)
			},
			expectedOwnerships: []*domain.Ownership{
				{ID: 1, ItemID: 1, UserID: "user1", Rentable: true, Memo: "memo1"},
				{ID: 2, ItemID: 1, UserID: "user2", Rentable: false, Memo: "memo2"},
			},
		},
		{
			name:   "failure",
			itemID: 2,
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByItemID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			tc.setupMock(mockOwnershipRepo)

			ownershipUseCase := NewOwnershipUseCase(mockOwnershipRepo)

			ownerships, err := ownershipUseCase.GetByItemID(tc.itemID)

			assert.Equal(t, tc.expectedOwnerships, ownerships)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestOwnershipUseCase_CreateOwnership(t *testing.T) {
	testCases := []struct {
		name              string
		inputOwnership    *domain.Ownership
		setupMock         func(repo *mock_domain.MockOwnershipRepository)
		expectedOwnership *domain.Ownership
		expectedErr       error
	}{
		{
			name: "success",
			inputOwnership: &domain.Ownership{
				ItemID:   1,
				UserID:   "user1",
				Rentable: true,
				Memo:     "new ownership",
			},
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					Create(&domain.Ownership{ItemID: 1, UserID: "user1", Rentable: true, Memo: "new ownership"}).
					Return(&domain.Ownership{ID: 1, ItemID: 1, UserID: "user1", Rentable: true, Memo: "new ownership"}, nil).
					Times(1)
			},
			expectedOwnership: &domain.Ownership{ID: 1, ItemID: 1, UserID: "user1", Rentable: true, Memo: "new ownership"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			tc.setupMock(mockOwnershipRepo)

			ownershipUseCase := NewOwnershipUseCase(mockOwnershipRepo)

			ownership, err := ownershipUseCase.CreateOwnership(tc.inputOwnership)

			assert.Equal(t, tc.expectedOwnership, ownership)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestOwnershipUseCase_UpdateOwnership(t *testing.T) {
	testCases := []struct {
		name              string
		inputOwnership    *domain.Ownership
		userID            string
		setupMock         func(repo *mock_domain.MockOwnershipRepository)
		expectedOwnership *domain.Ownership
		expectedErr       error
	}{
		{
			name: "success",
			inputOwnership: &domain.Ownership{
				ID:       1,
				ItemID:   1,
				UserID:   "owner",
				Rentable: false,
				Memo:     "updated memo",
			},
			userID: "owner",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: true, Memo: "before"}, nil).
					Times(1)

				repo.EXPECT().
					Update(&domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"}).
					Return(&domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"}, nil).
					Times(1)
			},
			expectedOwnership: &domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: false, Memo: "updated memo"},
		},
		{
			name: "failure: ownership not found",
			inputOwnership: &domain.Ownership{
				ID:       2,
				ItemID:   1,
				UserID:   "owner",
				Rentable: true,
				Memo:     "memo",
			},
			userID: "owner",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedErr: domain.ErrNotFound,
		},
		{
			name: "failure: cannot change other user's ownership",
			inputOwnership: &domain.Ownership{
				ID:       1,
				ItemID:   1,
				UserID:   "owner",
				Rentable: false,
				Memo:     "updated memo",
			},
			userID: "another-user",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: true, Memo: "before"}, nil).
					Times(1)
			},
			expectedErr: ErrForbidden,
		},
		{
			name: "failure: cannot change other user's ownership, although ownership.userID is matched",
			inputOwnership: &domain.Ownership{
				ID:       1,
				ItemID:   1,
				UserID:   "owner",
				Rentable: false,
				Memo:     "updated memo",
			},
			userID: "another-user",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: true, Memo: "before"}, nil).
					Times(1)
			},
			expectedErr: ErrForbidden,
		},
		{
			name: "failure: item mismatch",
			inputOwnership: &domain.Ownership{
				ID:       1,
				ItemID:   999,
				UserID:   "owner",
				Rentable: false,
				Memo:     "updated memo",
			},
			userID: "owner",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 1, UserID: "owner", Rentable: true, Memo: "before"}, nil).
					Times(1)
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			tc.setupMock(mockOwnershipRepo)

			ownershipUseCase := NewOwnershipUseCase(mockOwnershipRepo)

			ownership, err := ownershipUseCase.UpdateOwnership(tc.inputOwnership, tc.userID)

			assert.Equal(t, tc.expectedOwnership, ownership)
			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestOwnershipUseCase_DeleteOwnership(t *testing.T) {
	testCases := []struct {
		name        string
		id          int
		itemID      int
		userID      string
		setupMock   func(repo *mock_domain.MockOwnershipRepository)
		expectedErr error
	}{
		{
			name:   "success",
			id:     1,
			itemID: 10,
			userID: "owner",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 10, UserID: "owner", Rentable: true, Memo: "memo"}, nil).
					Times(1)

				repo.EXPECT().
					Delete(1).
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "failure: ownership not found",
			id:     2,
			itemID: 10,
			userID: "owner",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(2).
					Return(nil, domain.ErrNotFound).
					Times(1)
			},
			expectedErr: domain.ErrNotFound,
		},
		{
			name:   "failure: forbidden",
			id:     1,
			itemID: 10,
			userID: "another-user",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 10, UserID: "owner", Rentable: true, Memo: "memo"}, nil).
					Times(1)
			},
			expectedErr: ErrForbidden,
		},
		{
			name:   "failure: item mismatch",
			id:     1,
			itemID: 999,
			userID: "owner",
			setupMock: func(repo *mock_domain.MockOwnershipRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1, ItemID: 10, UserID: "owner", Rentable: true, Memo: "memo"}, nil).
					Times(1)
			},
			expectedErr: domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOwnershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			tc.setupMock(mockOwnershipRepo)

			ownershipUseCase := NewOwnershipUseCase(mockOwnershipRepo)

			err := ownershipUseCase.DeleteOwnership(tc.id, tc.itemID, tc.userID)

			assert.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}
