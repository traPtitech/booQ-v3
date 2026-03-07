package domain

import "time"

type Item struct {
	ID          int
	Name        string
	Description string
	ImgUrl      string

	BookDetail      *BookDetail
	EquipmentDetail *EquipmentDetail

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type BookDetail struct {
	ISBNCode string
}

type EquipmentDetail struct {
	Count    int
	CountMax int
}

type ItemSearchQuery struct {
	Name       string
	UserID     string
	BorrowerID string
	Limit      int
	Offset     int
	Tag        []string
	TagExclude []string
	SortBy     string
}

type ItemRepository interface {
	GetByID(id int) (*Item, error)
	Search(query ItemSearchQuery) ([]*Item, error)
	Create(item *Item) (*Item, error)
	CreateBatch(items []*Item) ([]*Item, error)
	Update(item *Item) (*Item, error)
	Delete(id int) error
}
