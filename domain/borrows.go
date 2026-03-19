package domain

import (
	"fmt"
	"time"
)

type BorrowingStatus string

const (
	BorrowingStatusRequested BorrowingStatus = "requested"
	BorrowingStatusBorrowed  BorrowingStatus = "borrowed"
	BorrowingStatusReturned  BorrowingStatus = "returned"
	BorrowingStatusRejected  BorrowingStatus = "rejected"
)

func (b BorrowingStatus) ToString() string {
	return string(b)
}

func ParseBorrowingStatus(s string) (BorrowingStatus, error) {
	switch s {
	case "requested":
		return BorrowingStatusRequested, nil
	case "borrowed":
		return BorrowingStatusBorrowed, nil
	case "returned":
		return BorrowingStatusReturned, nil
	case "rejected":
		return BorrowingStatusRejected, nil
	default:
		return "", fmt.Errorf("invalid borrowing status: %s", s)
	}
}

type Transaction struct {
	ID          int
	ItemID      int
	UserID      string // 借りる側
	OwnershipID int
	Status      BorrowingStatus

	// Request
	Purpose          string
	BorrowInClubRoom bool
	DueDate          time.Time // 返却予定日

	// Reply
	Message      string     // 貸す側のメッセージ
	CheckoutDate *time.Time // StatusがBorrowedになった日時

	// Return
	ReturnMessage string
	ReturnDate    *time.Time // StatusがReturnedになった日時

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTransaction(itemID int, userID string, ownershipID int, purpose string, borrowInClubRoom bool, dueDate time.Time) *Transaction {
	return &Transaction{
		ItemID:           itemID,
		UserID:           userID,
		OwnershipID:      ownershipID,
		Status:           BorrowingStatusRequested,
		Purpose:          purpose,
		BorrowInClubRoom: borrowInClubRoom,
		DueDate:          dueDate,
	}
}

func (t *Transaction) Approve(message string) error {
	if t.Status != BorrowingStatusRequested {
		return fmt.Errorf("%w: transaction is not in requested status", ErrInvalidTransactionStatus)
	}

	t.Status = BorrowingStatusBorrowed
	t.Message = message
	now := time.Now()
	t.CheckoutDate = &now
	return nil
}

func (t *Transaction) Reject(message string) error {
	if t.Status != BorrowingStatusRequested {
		return fmt.Errorf("%w: transaction is not in requested status", ErrInvalidTransactionStatus)
	}

	t.Status = BorrowingStatusRejected
	t.Message = message
	return nil
}

func (t *Transaction) Return(message string) error {
	if t.Status != BorrowingStatusBorrowed {
		return fmt.Errorf("%w: transaction is not in borrowed status", ErrInvalidTransactionStatus)
	}

	t.Status = BorrowingStatusReturned
	t.ReturnMessage = message
	now := time.Now()
	t.ReturnDate = &now
	return nil
}

type TransactionRepository interface {
	GetByID(id int) (*Transaction, error)
	GetByUserID(userID string) ([]*Transaction, error)
	GetByOwnershipID(ownershipID int) ([]*Transaction, error)
	Create(transaction *Transaction) (*Transaction, error)
	Update(transaction *Transaction) (*Transaction, error)
}
