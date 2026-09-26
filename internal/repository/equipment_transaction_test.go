package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/booQ-v3/internal/domain"
	"gorm.io/gorm"
)

var equipmentTestDueDate = time.Date(2030, 1, 2, 23, 59, 59, 0, time.UTC)

func createTestEquipmentTransaction(t *testing.T, repo domain.EquipmentTransactionRepository, userID string, itemID, count int) *domain.EquipmentTransaction {
	t.Helper()

	transaction, err := repo.Create(domain.NewEquipmentTransaction(userID, itemID, "purpose", count, equipmentTestDueDate))
	require.NoError(t, err)
	return transaction
}

func TestEquipmentTransactionRepository_Create(t *testing.T) {
	returnDate := time.Date(2030, 1, 3, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		name           string
		transaction    *domain.EquipmentTransaction
		wantStatus     int
		wantReturnDate bool
		wantErr        error
	}{
		{
			name:        "borrowed transaction",
			transaction: domain.NewEquipmentTransaction("user1", 1, "purpose", 3, equipmentTestDueDate),
			wantStatus:  1,
		},
		{
			name: "returned transaction",
			transaction: &domain.EquipmentTransaction{
				ItemID:        1,
				UserID:        "user1",
				Status:        domain.EquipmentBorrowingStatusReturned,
				Purpose:       "purpose",
				Count:         3,
				DueDate:       equipmentTestDueDate,
				ReturnMessage: "returned",
				ReturnDate:    returnDate,
			},
			wantStatus:     2,
			wantReturnDate: true,
		},
		{
			name:        "invalid status",
			transaction: &domain.EquipmentTransaction{Status: "invalid"},
			wantErr:     domain.ErrInvalidTransactionStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEquipmentTransactionRepository(db)

			created, err := repo.Create(tc.transaction)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, created)
				return
			}

			require.NoError(t, err)
			require.Positive(t, created.ID)
			require.False(t, created.CreatedAt.IsZero())
			require.False(t, created.UpdatedAt.IsZero())

			var stored equipmentTransaction
			require.NoError(t, db.First(&stored, created.ID).Error)
			assert.Equal(t, tc.wantStatus, stored.Status)
			assert.Equal(t, tc.transaction.UserID, stored.UserID)
			assert.Equal(t, tc.transaction.ItemID, stored.ItemID)
			assert.Equal(t, tc.transaction.Count, stored.Count)
			assert.Equal(t, tc.transaction.Purpose, stored.Purpose)
			assert.Equal(t, tc.transaction.ReturnMessage, stored.ReturnMessage)
			assert.True(t, equipmentTestDueDate.Equal(stored.DueDate))

			if tc.wantReturnDate {
				require.NotNil(t, stored.ReturnDate)
				assert.WithinDuration(t, returnDate, *stored.ReturnDate, time.Second)
			} else {
				assert.Nil(t, stored.ReturnDate)
			}
		})
	}
}

func TestEquipmentTransactionRepository_GetByID(t *testing.T) {
	testCases := []struct {
		name    string
		setup   func(t *testing.T, db *gorm.DB, repo domain.EquipmentTransactionRepository) int
		wantErr error
	}{
		{
			name: "found",
			setup: func(t *testing.T, _ *gorm.DB, repo domain.EquipmentTransactionRepository) int {
				return createTestEquipmentTransaction(t, repo, "user1", 1, 3).ID
			},
		},
		{
			name: "not found",
			setup: func(_ *testing.T, _ *gorm.DB, _ domain.EquipmentTransactionRepository) int {
				return 99999
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name: "invalid stored status",
			setup: func(t *testing.T, db *gorm.DB, _ domain.EquipmentTransactionRepository) int {
				model := &equipmentTransaction{
					ItemID:  1,
					UserID:  "user1",
					Purpose: "purpose",
					Count:   3,
					DueDate: equipmentTestDueDate,
					Status:  99,
				}
				require.NoError(t, db.Create(model).Error)
				return model.ID
			},
			wantErr: domain.ErrInvalidTransactionStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEquipmentTransactionRepository(db)
			id := tc.setup(t, db, repo)

			got, err := repo.GetByID(id)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, id, got.ID)
			assert.Equal(t, "user1", got.UserID)
			assert.Equal(t, 1, got.ItemID)
			assert.Equal(t, 3, got.Count)
			assert.Equal(t, "purpose", got.Purpose)
			assert.Equal(t, domain.EquipmentBorrowingStatusBorrowed, got.Status)
			assert.True(t, equipmentTestDueDate.Equal(got.DueDate))
			assert.True(t, got.ReturnDate.IsZero())
		})
	}
}

func TestEquipmentTransactionRepository_List(t *testing.T) {
	testCases := []struct {
		name        string
		find        func(repo domain.EquipmentTransactionRepository) ([]*domain.EquipmentTransaction, error)
		wantIndexes []int
	}{
		{
			name: "by user ID",
			find: func(repo domain.EquipmentTransactionRepository) ([]*domain.EquipmentTransaction, error) {
				return repo.GetByUserID("user1")
			},
			wantIndexes: []int{0, 2},
		},
		{
			name: "by item ID",
			find: func(repo domain.EquipmentTransactionRepository) ([]*domain.EquipmentTransaction, error) {
				return repo.GetByItemID(1)
			},
			wantIndexes: []int{0, 1},
		},
		{
			name: "unknown user ID",
			find: func(repo domain.EquipmentTransactionRepository) ([]*domain.EquipmentTransaction, error) {
				return repo.GetByUserID("missing")
			},
			wantIndexes: []int{},
		},
		{
			name: "unknown item ID",
			find: func(repo domain.EquipmentTransactionRepository) ([]*domain.EquipmentTransaction, error) {
				return repo.GetByItemID(99999)
			},
			wantIndexes: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEquipmentTransactionRepository(db)
			created := []*domain.EquipmentTransaction{
				createTestEquipmentTransaction(t, repo, "user1", 1, 3),
				createTestEquipmentTransaction(t, repo, "user2", 1, 1),
				createTestEquipmentTransaction(t, repo, "user1", 2, 2),
			}

			got, err := tc.find(repo)
			require.NoError(t, err)
			require.NotNil(t, got)

			wantIDs := make([]int, 0, len(tc.wantIndexes))
			for _, index := range tc.wantIndexes {
				wantIDs = append(wantIDs, created[index].ID)
			}

			gotIDs := make([]int, 0, len(got))
			for _, transaction := range got {
				gotIDs = append(gotIDs, transaction.ID)
			}
			assert.ElementsMatch(t, wantIDs, gotIDs)
		})
	}
}

func TestEquipmentTransactionRepository_Update(t *testing.T) {
	testCases := []struct {
		name             string
		initialReturned  bool
		prepare          func(transaction *domain.EquipmentTransaction) error
		wantStatus       domain.EquipmentBorrowingStatus
		wantStoredStatus int
		wantMessage      string
		wantReturnDate   bool
		wantErr          error
	}{
		{
			name: "return borrowed equipment",
			prepare: func(transaction *domain.EquipmentTransaction) error {
				return transaction.Return("all returned")
			},
			wantStatus:       domain.EquipmentBorrowingStatusReturned,
			wantStoredStatus: 2,
			wantMessage:      "all returned",
			wantReturnDate:   true,
		},
		{
			name:            "clear return fields",
			initialReturned: true,
			prepare: func(transaction *domain.EquipmentTransaction) error {
				transaction.Status = domain.EquipmentBorrowingStatusBorrowed
				transaction.ReturnMessage = ""
				transaction.ReturnDate = time.Time{}
				return nil
			},
			wantStatus:       domain.EquipmentBorrowingStatusBorrowed,
			wantStoredStatus: 1,
		},
		{
			name: "invalid status",
			prepare: func(transaction *domain.EquipmentTransaction) error {
				transaction.Status = "invalid"
				return nil
			},
			wantErr: domain.ErrInvalidTransactionStatus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewEquipmentTransactionRepository(db)
			created := createTestEquipmentTransaction(t, repo, "user1", 1, 3)

			if tc.initialReturned {
				require.NoError(t, created.Return("previous return"))
				_, err := repo.Update(created)
				require.NoError(t, err)
			}
			require.NoError(t, tc.prepare(created))

			updated, err := repo.Update(created)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, updated)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, created.ID, updated.ID)
			assert.False(t, updated.UpdatedAt.IsZero())

			got, err := repo.GetByID(created.ID)
			require.NoError(t, err)
			assert.Equal(t, tc.wantStatus, got.Status)
			assert.Equal(t, tc.wantMessage, got.ReturnMessage)
			assert.Equal(t, created.Count, got.Count)

			var stored equipmentTransaction
			require.NoError(t, db.First(&stored, created.ID).Error)
			assert.Equal(t, tc.wantStoredStatus, stored.Status)
			if tc.wantReturnDate {
				require.NotNil(t, stored.ReturnDate)
				assert.WithinDuration(t, created.ReturnDate, got.ReturnDate, time.Second)
			} else {
				assert.Nil(t, stored.ReturnDate)
				assert.True(t, got.ReturnDate.IsZero())
			}
		})
	}
}

func TestEquipmentTransactionInvalidStatus(t *testing.T) {
	testCases := []struct {
		name    string
		convert func() error
	}{
		{
			name: "domain to model",
			convert: func() error {
				_, err := toEquipmentTransactionModel(&domain.EquipmentTransaction{Status: "invalid"})
				return err
			},
		},
		{
			name: "model to domain",
			convert: func() error {
				_, err := (equipmentTransaction{Status: 99}).toDomain()
				return err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorIs(t, tc.convert(), domain.ErrInvalidTransactionStatus)
		})
	}
}
