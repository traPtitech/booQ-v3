package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

func createTestOwnership(t *testing.T, db *gorm.DB, userID string) *ownership {
	t.Helper()

	model := &ownership{
		ItemID:   1,
		UserID:   userID,
		Rentable: true,
		Memo:     "test ownership",
	}
	if err := db.Create(model).Error; err != nil {
		t.Fatalf("failed to create test ownership: %v", err)
	}

	return model
}

func TestTransactionRepository_GetByID(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) int
		expected *domain.Transaction
		wantErr  bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) int {
				o := createTestOwnership(t, db, "owner1")
				model := &transaction{
					UserID:           "user1",
					OwnershipID:      o.ID,
					Status:           string(domain.BorrowingStatusRequested),
					Purpose:          "test purpose",
					BorrowInClubRoom: false,
					DueDate:          time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				}
				if err := db.Create(model).Error; err != nil {
					t.Fatalf("failed to create test transaction: %v", err)
				}
				return model.ID
			},
			expected: &domain.Transaction{
				UserID:           "user1",
				OwnershipID:      1,
				Status:           domain.BorrowingStatusRequested,
				Purpose:          "test purpose",
				BorrowInClubRoom: false,
				DueDate:          time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "failure: transaction not found",
			setup: func(t *testing.T, db *gorm.DB) int {
				return 9999
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTransactionRepository(db)
			id := tc.setup(t, db)

			transaction, err := repo.GetByID(id)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.UserID, transaction.UserID)
				assert.Equal(t, tc.expected.OwnershipID, transaction.OwnershipID)
				assert.Equal(t, tc.expected.Status, transaction.Status)
				assert.Equal(t, tc.expected.Purpose, transaction.Purpose)
				assert.Equal(t, tc.expected.BorrowInClubRoom, transaction.BorrowInClubRoom)
				assert.True(t, tc.expected.DueDate.Equal(transaction.DueDate))
				assert.NotZero(t, transaction.ID)
				assert.NotZero(t, transaction.CreatedAt)
				assert.NotZero(t, transaction.UpdatedAt)
			}
		})
	}
}

func TestTransactionRepository_GetByUserID(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) string
		expected []*domain.Transaction
		wantErr  bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) string {
				o1 := createTestOwnership(t, db, "owner1")
				o2 := createTestOwnership(t, db, "owner2")
				o3 := createTestOwnership(t, db, "owner3")
				models := []*transaction{
					{
						UserID:      "user1",
						OwnershipID: o1.ID,
						Status:      "requested",
						Purpose:     "p1",
						DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						UserID:      "user1",
						OwnershipID: o2.ID,
						Status:      "requested",
						Purpose:     "p2",
						DueDate:     time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						UserID:      "user2",
						OwnershipID: o3.ID,
						Status:      "requested",
						Purpose:     "p3",
						DueDate:     time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC),
					},
				}
				for _, m := range models {
					db.Create(m)
				}
				return "user1"
			},
			expected: []*domain.Transaction{
				{
					UserID:      "user1",
					OwnershipID: 1,
					Status:      "requested",
					Purpose:     "p1",
					DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      "user1",
					OwnershipID: 2,
					Status:      "requested",
					Purpose:     "p2",
					DueDate:     time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "success: empty result",
			setup: func(t *testing.T, db *gorm.DB) string {
				return "nonexistent-user"
			},
			expected: []*domain.Transaction{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTransactionRepository(db)
			userID := tc.setup(t, db)

			transactions, err := repo.GetByUserID(userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, transactions)
			} else {
				assert.NoError(t, err)
				assert.Len(t, transactions, len(tc.expected))
				for i, tr := range transactions {
					assert.Equal(t, tc.expected[i].UserID, tr.UserID)
					assert.Equal(t, tc.expected[i].OwnershipID, tr.OwnershipID)
					assert.Equal(t, tc.expected[i].Status, tr.Status)
					assert.Equal(t, tc.expected[i].Purpose, tr.Purpose)
					assert.True(t, tc.expected[i].DueDate.Equal(tr.DueDate))
				}
			}
		})
	}
}

func TestTransactionRepository_GetByOwnershipID(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) int
		expected []*domain.Transaction
		wantErr  bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) int {
				o1 := createTestOwnership(t, db, "owner1")
				o2 := createTestOwnership(t, db, "owner2")
				models := []*transaction{
					{
						UserID:      "user1",
						OwnershipID: o1.ID,
						Status:      "requested",
						Purpose:     "p1",
						DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						UserID:      "user2",
						OwnershipID: o1.ID,
						Status:      "requested",
						Purpose:     "p2",
						DueDate:     time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC),
					},
					{
						UserID:      "user3",
						OwnershipID: o2.ID,
						Status:      "requested",
						Purpose:     "p3",
						DueDate:     time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC),
					},
				}
				for _, m := range models {
					db.Create(m)
				}
				return o1.ID
			},
			expected: []*domain.Transaction{
				{
					UserID:      "user1",
					OwnershipID: 1,
					Status:      "requested",
					Purpose:     "p1",
					DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					UserID:      "user2",
					OwnershipID: 1,
					Status:      "requested",
					Purpose:     "p2",
					DueDate:     time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "success: empty result",
			setup: func(t *testing.T, db *gorm.DB) int {
				return 9999
			},
			expected: []*domain.Transaction{},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTransactionRepository(db)
			ownershipID := tc.setup(t, db)

			transactions, err := repo.GetByOwnershipID(ownershipID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, transactions)
			} else {
				assert.NoError(t, err)
				assert.Len(t, transactions, len(tc.expected))
				for i, tr := range transactions {
					assert.Equal(t, tc.expected[i].UserID, tr.UserID)
					assert.Equal(t, tc.expected[i].OwnershipID, tr.OwnershipID)
					assert.Equal(t, tc.expected[i].Status, tr.Status)
					assert.Equal(t, tc.expected[i].Purpose, tr.Purpose)
					assert.True(t, tc.expected[i].DueDate.Equal(tr.DueDate))
				}
			}
		})
	}
}

func TestTransactionRepository_Create(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) *domain.Transaction
		expected *domain.Transaction
		wantErr  bool
	}{
		{
			name: "success",
			setup: func(t *testing.T, db *gorm.DB) *domain.Transaction {
				o := createTestOwnership(t, db, "owner1")
				return &domain.Transaction{
					UserID:           "user1",
					OwnershipID:      o.ID,
					Status:           domain.BorrowingStatusRequested,
					Purpose:          "create test",
					BorrowInClubRoom: true,
					DueDate:          time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				}
			},
			expected: &domain.Transaction{
				UserID:           "user1",
				OwnershipID:      1,
				Status:           domain.BorrowingStatusRequested,
				Purpose:          "create test",
				BorrowInClubRoom: true,
				DueDate:          time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "failure: ownership does not exist",
			setup: func(t *testing.T, db *gorm.DB) *domain.Transaction {
				return &domain.Transaction{
					UserID:      "user1",
					OwnershipID: 9999,
					Status:      domain.BorrowingStatusRequested,
					Purpose:     "create test",
					DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				}
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTransactionRepository(db)
			tr := tc.setup(t, db)

			created, err := repo.Create(tr)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, created)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.UserID, created.UserID)
				assert.Equal(t, tc.expected.OwnershipID, created.OwnershipID)
				assert.Equal(t, tc.expected.Status, created.Status)
				assert.Equal(t, tc.expected.Purpose, created.Purpose)
				assert.Equal(t, tc.expected.BorrowInClubRoom, created.BorrowInClubRoom)
				assert.True(t, tc.expected.DueDate.Equal(created.DueDate))
				assert.NotZero(t, created.ID)
				assert.NotZero(t, created.CreatedAt)
				assert.NotZero(t, created.UpdatedAt)

				var model transaction
				db.First(&model, created.ID)
				assert.Equal(t, tc.expected.UserID, model.UserID)
				assert.Equal(t, tc.expected.OwnershipID, model.OwnershipID)
				assert.Equal(t, string(tc.expected.Status), model.Status)
				assert.Equal(t, tc.expected.Purpose, model.Purpose)
				assert.Equal(t, tc.expected.BorrowInClubRoom, model.BorrowInClubRoom)
				assert.True(t, tc.expected.DueDate.Equal(model.DueDate))
			}
		})
	}
}

func TestTransactionRepository_Update(t *testing.T) {
	testCases := []struct {
		name     string
		setup    func(t *testing.T, db *gorm.DB) *domain.Transaction
		expected *domain.Transaction
		wantErr  bool
	}{
		{
			name: "success: approve transaction",
			setup: func(t *testing.T, db *gorm.DB) *domain.Transaction {
				o := createTestOwnership(t, db, "owner1")
				model := &transaction{
					UserID:      "user1",
					OwnershipID: o.ID,
					Status:      string(domain.BorrowingStatusRequested),
					Purpose:     "initial purpose",
					DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				}
				db.Create(model)

				tr := model.toDomain()
				tr.Approve("approved message")
				return tr
			},
			expected: &domain.Transaction{
				UserID:      "user1",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusBorrowed,
				Purpose:     "initial purpose",
				DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				Message:     "approved message",
			},
			wantErr: false,
		},
		{
			name: "failure: ownership does not exist",
			setup: func(t *testing.T, db *gorm.DB) *domain.Transaction {
				return &domain.Transaction{
					ID:          9999,
					UserID:      "user1",
					OwnershipID: 9999,
					Status:      domain.BorrowingStatusRequested,
					Purpose:     "test purpose",
					DueDate:     time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
				}
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "success: create new transaction when id does not exist",
			setup: func(t *testing.T, db *gorm.DB) *domain.Transaction {
				o := createTestOwnership(t, db, "owner2")
				return &domain.Transaction{
					ID:          9999,
					UserID:      "user2",
					OwnershipID: o.ID,
					Status:      domain.BorrowingStatusRequested,
					Purpose:     "new transaction",
					DueDate:     time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC),
				}
			},
			expected: &domain.Transaction{
				UserID:      "user2",
				OwnershipID: 1,
				Status:      domain.BorrowingStatusRequested,
				Purpose:     "new transaction",
				DueDate:     time.Date(2024, 7, 10, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := NewTransactionRepository(db)
			tr := tc.setup(t, db)

			updated, err := repo.Update(tr)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, updated)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.UserID, updated.UserID)
				assert.Equal(t, tc.expected.OwnershipID, updated.OwnershipID)
				assert.Equal(t, tc.expected.Status, updated.Status)
				assert.Equal(t, tc.expected.Purpose, updated.Purpose)
				assert.True(t, tc.expected.DueDate.Equal(updated.DueDate))
				assert.Equal(t, tc.expected.Message, updated.Message)
				assert.NotZero(t, updated.ID)
				assert.NotZero(t, updated.CreatedAt)
				assert.NotZero(t, updated.UpdatedAt)

				var model transaction
				db.First(&model, updated.ID)
				assert.Equal(t, tc.expected.UserID, model.UserID)
				assert.Equal(t, tc.expected.OwnershipID, model.OwnershipID)
				assert.Equal(t, string(tc.expected.Status), model.Status)
				assert.Equal(t, tc.expected.Purpose, model.Purpose)
				assert.True(t, tc.expected.DueDate.Equal(model.DueDate))
				assert.Equal(t, tc.expected.Message, model.Message)
			}
		})
	}
}
