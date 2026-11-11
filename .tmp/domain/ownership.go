package domain

type LendingPolicy string

const (
	Immediate        LendingPolicy = "immediate"
	ApprovalRequired LendingPolicy = "approval_required"
	Disable          LendingPolicy = "disable"
)

// D09: 独立した OwnershipID は設けない。
type OwnershipKey struct {
	ItemID    int
	OwnerKind OwnerKind
	OwnerID   string
}

func (k OwnershipKey) Owner() OwnerRef {
	return OwnerRef{Kind: k.OwnerKind, ID: k.OwnerID}
}

type Ownership struct {
	Key      OwnershipKey
	Quantity int // D03: 貸出中を含む所有総数。
	Policy   LendingPolicy
	Memo     string
}

func (o Ownership) Validate() error {
	if err := o.Key.Owner().Validate(); err != nil {
		return err
	}

	if o.Quantity < 0 {
		return ErrInvalidQuantity
	}

	switch o.Policy {
	case Immediate, ApprovalRequired, Disable:
		return nil
	}

	return ErrInvalidPolicy
}

type Availability struct {
	Total     int
	Borrowed  int
	Available int
}

func (o Ownership) Availability(borrowed int) Availability {
	return Availability{Total: o.Quantity, Borrowed: borrowed, Available: o.Quantity - borrowed}
}

func (o Ownership) RequirePolicy(policy LendingPolicy) error {
	if o.Policy == Disable {
		return ErrLendingDisabled
	}

	if o.Policy != policy {
		return ErrInvalidPolicy
	}

	return nil
}

// quantity: 借りたい数
func (o Ownership) CheckQuantity(quantity int, borrowed int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if quantity > o.Quantity-borrowed {
		return ErrInsufficientQuantity
	}

	return nil
}

func (o *Ownership) ChangeSettings(quantity int, policy LendingPolicy, memo string, borrowed int) error {
	next := *o
	next.Quantity, next.Policy, next.Memo = quantity, policy, memo
	if err := next.Validate(); err != nil {
		return err
	}

	if quantity < borrowed {
		return ErrInsufficientQuantity
	}

	*o = next
	return nil
}

type OwnershipRepository interface {
	Get(key OwnershipKey) (*Ownership, error)
	ListByItemID(itemID int) ([]*Ownership, error)
	Create(ownership *Ownership) (*Ownership, error)
	Save(ownership *Ownership) error
	Delete(key OwnershipKey) error
}
