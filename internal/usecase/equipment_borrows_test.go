package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/booQ-v3/internal/domain"
	mock_domain "github.com/traPtitech/booQ-v3/internal/domain/mock"
	"go.uber.org/mock/gomock"
)

func TestBorrowEquipment(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	for _, tc := range []struct {
		name                       string
		count                      int
		due                        time.Time
		item                       *domain.Item
		getErr, createErr, wantErr error
	}{
		{
			name:  "success",
			count: 3,
			due:   future,
			item:  &domain.Item{EquipmentDetail: &domain.EquipmentDetail{}},
		},
		{
			name:  "default count",
			count: 0,
			due:   future,
			item:  &domain.Item{EquipmentDetail: &domain.EquipmentDetail{}},
		},
		{
			name:  "negative count",
			count: -1,
			due:   future,
			item:  &domain.Item{EquipmentDetail: &domain.EquipmentDetail{}},
		},
		{
			name:    "missing item",
			getErr:  domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
		{
			name:    "item read error",
			getErr:  assert.AnError,
			wantErr: assert.AnError,
		},
		{
			name:    "book",
			item:    &domain.Item{},
			wantErr: ErrItemNotEquipment,
		},
		{
			name:    "past date",
			item:    &domain.Item{EquipmentDetail: &domain.EquipmentDetail{}},
			due:     time.Now().Add(-time.Hour),
			wantErr: ErrInvalidDueDate,
		},
		{
			name:    "missing date",
			item:    &domain.Item{EquipmentDetail: &domain.EquipmentDetail{}},
			wantErr: ErrInvalidDueDate,
		},
		{
			name:      "create error",
			count:     2,
			due:       future,
			item:      &domain.Item{EquipmentDetail: &domain.EquipmentDetail{}},
			createErr: assert.AnError,
			wantErr:   assert.AnError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			items := mock_domain.NewMockItemRepository(ctrl)
			transactions := mock_domain.NewMockEquipmentTransactionRepository(ctrl)

			items.EXPECT().GetByID(7).Return(tc.item, tc.getErr)
			if tc.wantErr == nil || tc.createErr != nil {
				count := tc.count
				if count <= 0 {
					count = 1
				}
				transactions.EXPECT().Create(domain.NewEquipmentTransaction("user", 7, "purpose", count, tc.due)).DoAndReturn(func(tr *domain.EquipmentTransaction) (*domain.EquipmentTransaction, error) {
					if tc.createErr != nil {
						return nil, tc.createErr
					}
					tr.ID = 12
					return tr, nil
				})
			}

			u := NewBorrowingUseCase(nil, nil, transactions, items)
			got, err := u.BorrowEquipment(7, "user", "purpose", tc.count, tc.due)

			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr != nil {
				require.Nil(t, got)
			} else {
				assert.Equal(t, 12, got.ID)
				assert.True(t, got.ReturnDate.IsZero())
			}
		})
	}
}

func TestReturnEquipment(t *testing.T) {
	for _, tc := range []struct {
		name, message              string
		itemID                     int
		userID                     string
		status                     domain.EquipmentBorrowingStatus
		getErr, updateErr, wantErr error
	}{
		{
			name:    "success",
			message: "returned",
			itemID:  7,
			userID:  "user",
			status:  domain.EquipmentBorrowingStatusBorrowed,
		},
		{
			name:   "empty message",
			itemID: 7,
			userID: "user",
			status: domain.EquipmentBorrowingStatusBorrowed,
		},
		{
			name:    "not found",
			getErr:  domain.ErrNotFound,
			wantErr: domain.ErrNotFound,
		},
		{
			name:    "different item",
			itemID:  8,
			userID:  "user",
			status:  domain.EquipmentBorrowingStatusBorrowed,
			wantErr: domain.ErrNotFound,
		},
		{
			name:    "different user",
			itemID:  7,
			userID:  "other",
			status:  domain.EquipmentBorrowingStatusBorrowed,
			wantErr: ErrForbidden,
		},
		{
			name:    "already returned",
			itemID:  7,
			userID:  "user",
			status:  domain.EquipmentBorrowingStatusReturned,
			wantErr: domain.ErrInvalidTransactionStatus,
		},
		{
			name:    "read error",
			getErr:  assert.AnError,
			wantErr: assert.AnError,
		},
		{
			name:      "update error",
			itemID:    7,
			userID:    "user",
			status:    domain.EquipmentBorrowingStatusBorrowed,
			updateErr: assert.AnError,
			wantErr:   assert.AnError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := mock_domain.NewMockEquipmentTransactionRepository(gomock.NewController(t))

			var transaction *domain.EquipmentTransaction
			if tc.getErr == nil {
				transaction = &domain.EquipmentTransaction{
					ID:     4,
					ItemID: tc.itemID,
					UserID: tc.userID,
					Status: tc.status,
					Count:  3,
				}
			}

			repo.EXPECT().GetByID(4).Return(transaction, tc.getErr)
			if tc.wantErr == nil || tc.updateErr != nil {
				repo.EXPECT().Update(gomock.Any()).DoAndReturn(func(tr *domain.EquipmentTransaction) (*domain.EquipmentTransaction, error) {
					assert.Equal(t, 4, tr.ID)
					assert.Equal(t, 3, tr.Count)
					assert.Equal(t, domain.EquipmentBorrowingStatusReturned, tr.Status)
					assert.Equal(t, tc.message, tr.ReturnMessage)
					assert.WithinDuration(t, time.Now(), tr.ReturnDate, time.Second)
					if tc.updateErr != nil {
						return nil, tc.updateErr
					}
					return tr, nil
				})
			}

			got, err := NewBorrowingUseCase(nil, nil, repo, nil).ReturnEquipment(7, 4, "user", tc.message)

			require.ErrorIs(t, err, tc.wantErr)
			if tc.wantErr != nil {
				assert.Nil(t, got)
				if transaction != nil && tc.updateErr == nil {
					assert.Equal(t, tc.status, transaction.Status)
					assert.True(t, transaction.ReturnDate.IsZero())
				}
			} else {
				require.NotNil(t, got)
				assert.Equal(t, 4, got.ID)
			}
		})
	}
}
