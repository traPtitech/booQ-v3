package domain

import "time"

type TransactionStatus string

const (
	Borrowed TransactionStatus = "borrowed"
	Returned TransactionStatus = "returned"
)

type Transaction struct {
	ID               int
	OwnershipKey     OwnershipKey
	BorrowerID       string
	Quantity         int
	Purpose          string
	DueAt            time.Time
	BorrowInClubRoom bool
	RequestID        *int // 自由貸出では nil。
	Status           TransactionStatus
	CheckedOutAt     time.Time
	ReturnedAt       *time.Time
	ReturnMessage    string
}

func NewTransaction(
	ownershipKey OwnershipKey,
	borrowerID string,
	quantity int,
	purpose string,
	dueAt time.Time,
	borrowInClubRoom bool,
	requestID *int,
	now time.Time,
) (*Transaction, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	if !dueAt.After(now) {
		return nil, ErrInvalidDueDate
	}

	return &Transaction{
		OwnershipKey:     ownershipKey,
		BorrowerID:       borrowerID,
		Quantity:         quantity,
		Purpose:          purpose,
		DueAt:            dueAt,
		BorrowInClubRoom: borrowInClubRoom,
		RequestID:        requestID,
		Status:           Borrowed,
		CheckedOutAt:     now,
	}, nil
}

func (t *Transaction) Return(message string, now time.Time) error {
	if t.Status != Borrowed {
		return ErrInvalidState
	}

	// D10: 全量返却。期限を超えていても返却可能。
	t.Status, t.ReturnMessage, t.ReturnedAt = Returned, message, &now
	return nil
}

type TransactionRepository interface {
	GetByID(id int) (*Transaction, error)
	ListByOwnership(key OwnershipKey) ([]*Transaction, error)
	SumBorrowedQuantity(key OwnershipKey) (int, error)
	Create(transaction *Transaction) (*Transaction, error)
	Save(transaction *Transaction) error
}
