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

func NewTransaction(userID string, ownershipID int, purpose string, borrowInClubRoom bool, dueDate time.Time) *Transaction {
	return &Transaction{
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

type EquipmentBorrowingStatus string

const (
	EquipmentBorrowingStatusBorrowed EquipmentBorrowingStatus = "borrowed"
	EquipmentBorrowingStatusReturned EquipmentBorrowingStatus = "returned"
)

type EquipmentTransaction struct {
	ID      int
	ItemID  int
	UserID  string
	Status  EquipmentBorrowingStatus
	Purpose string
	Count   int
	DueDate time.Time

	ReturnMessage string
	ReturnDate    time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEquipmentTransaction(userID string, itemID int, purpose string, count int, dueDate time.Time) *EquipmentTransaction {
	return &EquipmentTransaction{
		UserID:  userID,
		ItemID:  itemID,
		Status:  EquipmentBorrowingStatusBorrowed,
		Purpose: purpose,
		Count:   count,
		DueDate: dueDate,
	}
}

func (t *EquipmentTransaction) Return(message string) error {
	if t.Status != EquipmentBorrowingStatusBorrowed {
		return fmt.Errorf("%w: transaction is not in borrowed status", ErrInvalidTransactionStatus)
	}

	t.Status = EquipmentBorrowingStatusReturned
	t.ReturnMessage = message
	t.ReturnDate = time.Now()
	return nil
}

type EquipmentTransactionRepository interface {
	GetByID(id int) (*EquipmentTransaction, error)
	GetByUserID(userID string) ([]*EquipmentTransaction, error)
	GetByItemID(itemID int) ([]*EquipmentTransaction, error)
	Create(transaction *EquipmentTransaction) (*EquipmentTransaction, error)
	Update(transaction *EquipmentTransaction) (*EquipmentTransaction, error)	
}
