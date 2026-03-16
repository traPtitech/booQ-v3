package domain

import "time"

type BorrowingStatus int

const (
	BorrowingStatusBorrowed BorrowingStatus = 1
	BorrowingStatusReturned BorrowingStatus = 2
)

type Transaction struct {
	ID            int
	ItemID        int
	UserID        string
	OwnershipID   int
	Status        BorrowingStatus
	Purpose       string
	Message       string
	ReturnMessage string
	DueDate       time.Time
	CheckoutDate  time.Time
	ReturnDate    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

type TransactionRepository interface {
	GetByID(id int) (*Transaction, error)
	GetByUserID(userID string) ([]*Transaction, error)
	GetByOwnershipID(ownershipID int) ([]*Transaction, error)
	Create(transaction *Transaction) (*Transaction, error)
	Update(transaction *Transaction) (*Transaction, error)
}
