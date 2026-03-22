package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBorrowingStatus_ToString(t *testing.T) {
	tests := []struct {
		name     string
		s        BorrowingStatus
		expected string
	}{
		{
			name:     "requested",
			s:        BorrowingStatusRequested,
			expected: "requested",
		},
		{
			name:     "borrowed",
			s:        BorrowingStatusBorrowed,
			expected: "borrowed",
		},
		{
			name:     "returned",
			s:        BorrowingStatusReturned,
			expected: "returned",
		},
		{
			name:     "rejected",
			s:        BorrowingStatusRejected,
			expected: "rejected",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.s.ToString())
		})
	}
}

func TestParseBorrowingStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BorrowingStatus
		wantErr  bool
	}{
		{
			name:     "requested",
			input:    "requested",
			expected: BorrowingStatusRequested,
			wantErr:  false,
		},
		{
			name:     "borrowed",
			input:    "borrowed",
			expected: BorrowingStatusBorrowed,
			wantErr:  false,
		},
		{
			name:     "returned",
			input:    "returned",
			expected: BorrowingStatusReturned,
			wantErr:  false,
		},
		{
			name:     "rejected",
			input:    "rejected",
			expected: BorrowingStatusRejected,
			wantErr:  false,
		},
		{
			name:    "invalid",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBorrowingStatus(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestNewTransaction(t *testing.T) {
	dueDate := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name             string
		itemID           int
		userID           string
		ownershipID      int
		purpose          string
		borrowInClubRoom bool
		dueDate          time.Time
	}{
		{
			name:             "success",
			itemID:           1,
			userID:           "user1",
			ownershipID:      1,
			purpose:          "purpose",
			borrowInClubRoom: true,
			dueDate:          dueDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTransaction(tt.userID, tt.ownershipID, tt.purpose, tt.borrowInClubRoom, tt.dueDate)
			assert.Equal(t, tt.userID, tr.UserID)
			assert.Equal(t, tt.ownershipID, tr.OwnershipID)
			assert.Equal(t, tt.purpose, tr.Purpose)
			assert.Equal(t, tt.borrowInClubRoom, tr.BorrowInClubRoom)
			assert.Equal(t, tt.dueDate, tr.DueDate)
			assert.Equal(t, BorrowingStatusRequested, tr.Status)
		})
	}
}

func TestTransaction_Approve(t *testing.T) {
	tests := []struct {
		name    string
		status  BorrowingStatus
		message string
		wantErr bool
	}{
		{
			name:    "success",
			status:  BorrowingStatusRequested,
			message: "ok",
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  BorrowingStatusBorrowed,
			message: "ok",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transaction{Status: tt.status}
			err := tr.Approve(tt.message)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalidTransactionStatus)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, BorrowingStatusBorrowed, tr.Status)
				assert.Equal(t, tt.message, tr.Message)
				assert.NotNil(t, tr.CheckoutDate)
			}
		})
	}
}

func TestTransaction_Reject(t *testing.T) {
	tests := []struct {
		name    string
		status  BorrowingStatus
		message string
		wantErr bool
	}{
		{
			name:    "success",
			status:  BorrowingStatusRequested,
			message: "no",
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  BorrowingStatusBorrowed,
			message: "no",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transaction{Status: tt.status}
			err := tr.Reject(tt.message)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalidTransactionStatus)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, BorrowingStatusRejected, tr.Status)
				assert.Equal(t, tt.message, tr.Message)
			}
		})
	}
}

func TestTransaction_Return(t *testing.T) {
	tests := []struct {
		name    string
		status  BorrowingStatus
		message string
		wantErr bool
	}{
		{
			name:    "success",
			status:  BorrowingStatusBorrowed,
			message: "returned",
			wantErr: false,
		},
		{
			name:    "invalid status",
			status:  BorrowingStatusRequested,
			message: "returned",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transaction{Status: tt.status}
			err := tr.Return(tt.message)
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalidTransactionStatus)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, BorrowingStatusReturned, tr.Status)
				assert.Equal(t, tt.message, tr.ReturnMessage)
				assert.NotNil(t, tr.ReturnDate)
			}
		})
	}
}
