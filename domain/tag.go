package domain

type Tag struct {
	ItemID int
	Name   string
}

type TagRepository interface {
	GetByItemID(itemID int) ([]*Tag, error)
	ReplaceByItemID(itemID int, tags []string) error
}
