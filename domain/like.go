package domain

type Like struct {
	ItemID int
	UserID string
}

type LikeRepository interface {
	GetByItemID(itemID int) ([]*Like, error)
	Exists(itemID int, userID string) (bool, error)
	Create(like *Like) error
	Delete(itemID int, userID string) error
}
