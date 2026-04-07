package usecase

import (
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
		userID              string
		ownershipID         int
		purpose             string
		dueDate             time.Time
		borrowInClubRoom    bool
		setupMock           func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository)
		expectedTransaction *domain.Transaction
		expectedError       error
	}{
		{
			name:             "success",
			userID:           "user1",
			ownershipID:      1,
			purpose:          "for study",
			dueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{ID: 1}, nil)
				transactionRepo.EXPECT().
					Create(&domain.Transaction{
						UserID:           "user1",
						OwnershipID:      1,
						Status:           domain.BorrowingStatusRequested,
						Purpose:          "for study",
						DueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
						BorrowInClubRoom: false,
					}).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
						Purpose:     "for study",
						DueDate:     time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
						Status:      domain.BorrowingStatusRequested,
					}, nil)
			},
			expectedTransaction: &domain.Transaction{
				ID:          1,
				UserID:      "user1",
				OwnershipID: 1,
				Purpose:     "for study",
				DueDate:     time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
				Status:      domain.BorrowingStatusRequested,
			},
			expectedError: nil,
		},
		{
			name:             "failure: ownership not found",
			userID:           "user1",
			ownershipID:      999,
			purpose:          "for study",
			dueDate:          time.Date(2200, 7, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				ownershipRepo.EXPECT().
					GetByID(999).
					Return(nil, domain.ErrNotFound)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
		{
			name:             "failure: due date in the past",
			userID:           "user1",
			ownershipID:      1,
			purpose:          "for study",
			dueDate:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			borrowInClubRoom: false,
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
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

			ownershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(ownershipRepo, transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, ownershipRepo)
			transaction, err := u.PostRequest(tc.userID, tc.ownershipID, tc.purpose, tc.dueDate, tc.borrowInClubRoom)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.Equal(t, tc.expectedTransaction, transaction)
		})
	}
}

func TestBorrowingUseCase_GetRequest(t *testing.T) {
	testCases := []struct {
		name                string
		userID              string
		ownershipID         int
		borrowingID         int
		setupMock           func(transactionRepo *mock_domain.MockTransactionRepository)
		expectedTransaction *domain.Transaction
		expectedError       error
	}{
		{
			name:        "success",
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
			},
			expectedTransaction: &domain.Transaction{
				ID:          1,
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusRequested,
			},
			expectedError: nil,
		},
		{
			name:        "failure: transaction not found",
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
			name:        "failure: userID mismatch",
			userID:      "user2",
			ownershipID: 1,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       ErrForbidden,
		},
		{
			name:        "failure: ownershipID mismatch",
			userID:      "user1",
			ownershipID: 2,
			borrowingID: 1,
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, nil)
			transaction, err := u.GetRequest(tc.userID, tc.ownershipID, tc.borrowingID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.Equal(t, tc.expectedTransaction, transaction)
		})
	}
}

func TestBorrowingUseCase_ReplyRequest(t *testing.T) {
	testCases := []struct {
		name                string
		userID              string
		ownershipID         int
		borrowingID         int
		approve             bool
		message             string
		setupMock           func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository)
		expectedTransaction *domain.Transaction
		expectedError       error
	}{
		{
			name:        "success: approve",
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{
						ID:     1,
						UserID: "user1",
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
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusBorrowed,
				Message:     "approved",
			},
			expectedError: nil,
		},
		{
			name:        "success: reject",
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     false,
			message:     "rejected",
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{
						ID:     1,
						UserID: "user1",
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
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusRejected,
				Message:     "rejected",
			},
			expectedError: nil,
		},
		{
			name:        "failure: get request error",
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(nil, domain.ErrNotFound)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrNotFound,
		},
		{
			name:        "failure: invalid status",
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "user1",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusBorrowed,
					}, nil)
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{
						ID:     1,
						UserID: "user1",
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       domain.ErrInvalidTransactionStatus,
		},
		{
			name:        "failure: ownership owner mismatch",
			userID:      "user2",
			ownershipID: 1,
			borrowingID: 1,
			approve:     true,
			message:     "approved",
			setupMock: func(ownershipRepo *mock_domain.MockOwnershipRepository, transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
						UserID:      "borrower",
						OwnershipID: 1,
						Status:      domain.BorrowingStatusRequested,
					}, nil)
				ownershipRepo.EXPECT().
					GetByID(1).
					Return(&domain.Ownership{
						ID:     1,
						UserID: "user1",
					}, nil)
			},
			expectedTransaction: nil,
			expectedError:       ErrForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ownershipRepo := mock_domain.NewMockOwnershipRepository(ctrl)
			transactionRepo := mock_domain.NewMockTransactionRepository(ctrl)

			tc.setupMock(ownershipRepo, transactionRepo)

			u := NewBorrowingUseCase(transactionRepo, ownershipRepo)
			transaction, err := u.ReplyRequest(tc.userID, tc.ownershipID, tc.borrowingID, tc.approve, tc.message)

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
		userID        string
		ownershipID   int
		borrowingID   int
		message       string
		setupMock     func(transactionRepo *mock_domain.MockTransactionRepository)
		expectedError error
	}{
		{
			name:        "success",
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			message:     "returned",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
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
			userID:      "user1",
			ownershipID: 1,
			borrowingID: 1,
			message:     "returned",
			setupMock: func(transactionRepo *mock_domain.MockTransactionRepository) {
				transactionRepo.EXPECT().
					GetByID(1).
					Return(&domain.Transaction{
						ID:          1,
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

			u := NewBorrowingUseCase(transactionRepo, nil)
			err := u.ReturnItem(tc.userID, tc.ownershipID, tc.borrowingID, tc.message)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
