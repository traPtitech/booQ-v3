package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	mock_domain "github.com/traPtitech/booQ-v3/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestBorrowingUseCase_PostRequest(t *testing.T) {
	testCases := []struct {
		name                string
		itemID              int
		userID              string
		ownershipID         int
		purpose             string
		dueDate             time.Time
		borrowInClubRoom    bool
		setupMock           func(itemRepo *mock_domain.MockItemRepository, ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository)
		expectedTransaction *domain.Transaction
		expectedError       error
	}{
		{
			name:             "success",
			itemID:           1,
			userID:           "user1",
			ownershipID:      1,
			purpose:          "for study",
			dueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(itemRepo *mock_domain.MockItemRepository, ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				itemRepo.EXPECT().
					GetByID(1).
					Return(&domain.Item{ID: 1}, nil)
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1}, nil)
				transactionRepo.EXPECT().
					Create(&domain.Transaction{
						ItemID:           1,
						UserID:           "user1",
						OwnershipID:      1,
						Status:           domain.BorrowingStatusRequested,
						Purpose:          "for study",
						DueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
						BorrowInClubRoom: false,
					}).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Purpose:     "for study",
						DueDate:     time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
						Status:      domain.BorrowingStatusRequested,
					}, nil)
			},
			expectedTransaction: &domain.Transaction{
				ID:          1,
				ItemID:      1,
				UserID:      "user1",
				OwnershipID: 1,
				Purpose:     "for study",
				DueDate:     time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
				Status:      domain.BorrowingStatusRequested,
			},
			expectedError: nil,
		},
		{
			name:             "failure: item not found",
			itemID:           999,
			userID:           "user1",
			ownershipID:      1,
			purpose:          "for study",
			dueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(itemRepo *mock_domain.MockItemRepository, ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				itemRepo.EXPECT().
					GetByID(999).
					Return(nil, domain.ErrNotFound)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
		{
			name:             "failure: ownership not found",
			itemID:           1,
			userID:           "user1",
			ownershipID:      999,
			purpose:          "for study",
			dueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(itemRepo *mock_domain.MockItemRepository, ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				itemRepo.EXPECT().
					GetByID(1).
					Return(&domain.Item{ID: 1}, nil)
				ownershipRepo.EXPECT().
					GetByID(999).
					Return(nil, domain.ErrNotFound)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
		{
			name:             "failure: due date in the past",
			itemID:           1,
			userID:           "user1",
			ownershipID:      1,
			purpose:          "for study",
			dueDate:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(itemRepo *mock_domain.MockItemRepository, ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				itemRepo.EXPECT().
					GetByID(1).
					Return(&domain.Item{ID: 1}, nil)
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1}, nil)
			},
			expectedTransaction: nil,
			expectedError:       ErrInvalidDueDate,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			itemRepo := mock_domain.NewMockItemRepository(ctrl)
			ownershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(itemRepo, ownershipRepo, transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, itemRepo, ownershipRepo)
			transaction, err := u.PostRequest(tc.itemID, tc.userID, tc.ownershipID, tc.purpose, tc.dueDate, tc.borrowInClubRoom)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert .Equal(t, tc.expectedTransaction, transaction)
		})
	}
}

func TestBorrowingUseCase_GetRequest(t *testing.T) {
	testCases := []struct {
		name                string
		itemID              int
		userID              string
		ownershipID         int
		borrowingID         int
		setupMock           func(transactionRepo *mock_domain.MockTransactionRepository)
		expectedTransaction *domain.Transaction
		expectedError       error
	}{
		{
			name:        "success",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
			},
			expectedTransaction: &domain.Transaction{
				ID:          1,
				ItemID:      1,
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusRequested,
			},
			expectedError: nil,
		},
		{
			name:        "failure: transaction not found",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 999,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(999).
					Return(nil, domain.ErrNotFound)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
		{
			name:        "failure: itemID mismatch",
			itemID:      2,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       assert.AnError, // Will check error message or use a more specific error if possible
		},
		{
			name:        "failure: userID mismatch",
			itemID:      1,
			userID:      "user2",
			ownershipID: 1,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       assert.AnError,
		},
		{
			name:        "failure: ownershipID mismatch",
			itemID:      1,
			userID:      "user1",
			ownershipID: 2,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, nil, nil)
			transaction, err := u.GetRequest(tc.itemID, tc.userID, tc.ownershipID, tc.borrowingID)

			if tc.expectedError != nil {
				if errors.Is(tc.expectedError, assert.AnError) {
					assert.Error(t, err)
				} else {
					assert.ErrorIs(t, err, tc.expectedError)
				}
				return
			}

			assert.Equal(t, tc.expectedTransaction, transaction)
		})
	}
}

func TestBorrowingUseCase_ReplyRequest(t *testing.T) {
	testCases := []struct {
		name                string
		itemID              int
		userID              string
		ownershipID         int
		borrowingID         int
		approve             bool
		message             string
		setupMock           func(transactionRepo *mock_domain.MockTransactionRepository)
		expectedTransaction *domain.Transaction
		expectedError       error
	}{
		{
			name:        "success: approve",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
				transactionRepo.EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(tr *domain.Transaction) (*domain.Transaction, error) {
						assert.Equal(t, domain.BorrowingStatusBorrowed, tr.Status)
						assert.Equal(t, "approved", tr.Message)
						return tr, nil
					})
			},
			expectedTransaction: &domain.Transaction{
				ID:          1,
				ItemID:      1,
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusBorrowed,
				Message:     "approved",
			},
			expectedError: nil,
		},
		{
			name:        "success: reject",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     false,
			message:     "rejected",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
				transactionRepo.EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(tr *domain.Transaction) (*domain.Transaction, error) {
						assert.Equal(t, domain.BorrowingStatusRejected, tr.Status)
						assert.Equal(t, "rejected", tr.Message)
						return tr, nil
					})
			},
			expectedTransaction: &domain.Transaction{
				ID:          1,
				ItemID:      1,
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusRejected,
				Message:     "rejected",
			},
			expectedError: nil,
		},
		{
			name:        "failure: get request error",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(nil, domain.ErrNotFound)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
		{
			name:        "failure: invalid status",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusBorrowed,
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrInvalidTransactionStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, nil, nil)
			transaction, err := u.ReplyRequest(tc.itemID, tc.userID, tc.ownershipID, tc.borrowingID, tc.approve, tc.message)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.Equal(t, tc.expectedTransaction.ID, transaction.ID)
			assert.Equal(t, tc.expectedTransaction.Status, transaction.Status)
			assert.Equal(t, tc.expectedTransaction.Message, transaction.Message)
		})
	}
}

func TestBorrowingUseCase_ReturnItem(t *testing.T) {
	testCases := []struct {
		name          string
		itemID        int
		userID        string
		ownershipID   int
		borrowingID   int
		message       string
		setupMock     func(transactionRepo *mock_domain.MockTransactionRepository)
		expectedError error
	}{
		{
			name:        "success",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			message:     "returned",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusBorrowed,
					}, nil)
				transactionRepo.EXPECT().
					Update(gomock.Any()).
					DoAndReturn(func(tr *domain.Transaction) (*domain.Transaction, error) {
						assert.Equal(t, domain.BorrowingStatusReturned, tr.Status)
						assert.Equal(t, "returned", tr.ReturnMessage)
						return tr, nil
					})
			},
			expectedError: nil,
		},
		{
			name:        "failure: get request error",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			message:     "returned",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(nil, domain.ErrNotFound)
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name:        "failure: invalid status",
			itemID:      1,
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			message:     "returned",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						ItemID:      1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
			},
			expectedError: domain.ErrInvalidTransactionStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, nil, nil)
			err := u.ReturnItem(tc.itemID, tc.userID, tc.ownershipID, tc.borrowingID, tc.message)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
