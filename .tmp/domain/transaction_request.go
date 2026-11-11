package domain

import "time"

type RequestStatus string

const (
	Pending   RequestStatus = "pending"
	Approved  RequestStatus = "approved"
	Rejected  RequestStatus = "rejected"
	Cancelled RequestStatus = "cancelled"
)

type TransactionRequest struct {
	ID               int
	OwnershipKey     OwnershipKey
	BorrowerID       string
	Quantity         int
	Purpose          string
	DueAt            time.Time
	BorrowInClubRoom bool
	Status           RequestStatus
	CreatedAt        time.Time
	DecidedAt        *time.Time
	DecidedBy        string // 許可した人 所有者 or Groupなら管理者
	ReplyMessage     string
}

func NewRequest(
	ownershipKey OwnershipKey,
	borrowerID string,
	quantity int,
	purpose string,
	dueAt time.Time,
	borrowInClubRoom bool,
	now time.Time,
) (*TransactionRequest, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	if !dueAt.After(now) {
		return nil, ErrInvalidDueDate
	}

	return &TransactionRequest{
		OwnershipKey:     ownershipKey,
		BorrowerID:       borrowerID,
		Quantity:         quantity,
		Purpose:          purpose,
		DueAt:            dueAt,
		BorrowInClubRoom: borrowInClubRoom,
		Status:           Pending,
		CreatedAt:        now,
	}, nil
}

func (r *TransactionRequest) decide(status RequestStatus, actor string, message string, now time.Time) error {
	if r.Status != Pending {
		return ErrInvalidState
	}

	r.Status, r.DecidedBy, r.ReplyMessage, r.DecidedAt = status, actor, message, &now
	return nil
}

func (r *TransactionRequest) Approve(actor string, message string, now time.Time) (*Transaction, error) {
	if r.Status != Pending {
		return nil, ErrInvalidState
	}

	// D06: 承認で即貸出。申請内容を値として貸出へコピーする。
	loan, err := NewTransaction(
		r.OwnershipKey,
		r.BorrowerID,
		r.Quantity,
		r.Purpose,
		r.DueAt,
		r.BorrowInClubRoom,
		&r.ID,
		now,
	)
	if err != nil {
		return nil, err
	}

	if err := r.decide(Approved, actor, message, now); err != nil {
		return nil, err
	}

	return loan, nil
}

func (r *TransactionRequest) Reject(actor string, message string, now time.Time) error {
	return r.decide(Rejected, actor, message, now)
}

func (r *TransactionRequest) Cancel(actor string, now time.Time) error {
	return r.decide(Cancelled, actor, "", now)
}

type RequestRepository interface {
	GetByID(id int) (*TransactionRequest, error)
	ListByOwnership(key OwnershipKey) ([]*TransactionRequest, error)
	HasPending(key OwnershipKey) (bool, error)
	Create(request *TransactionRequest) (*TransactionRequest, error)
	Save(request *TransactionRequest) error
	SaveApproval(request *TransactionRequest, transaction *Transaction) (*Transaction, error)
}
