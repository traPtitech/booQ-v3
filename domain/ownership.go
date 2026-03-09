package domain

type Ownership struct {
	ID       int
	ItemID   int
	UserID   string
	Rentable bool
	Memo     string
}

type OwnershipRepository interface {
	GetByID(id int) (*Ownership, error)
	GetByItemID(itemID int) ([]*Ownership, error)
	GetByUserID(userID string) ([]*Ownership, error)
	Create(ownership *Ownership) (*Ownership, error)
	Update(ownership *Ownership) (*Ownership, error)
	Delete(id int) error
}
