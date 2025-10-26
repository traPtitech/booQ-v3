package domain

type Item struct {
	ID          int
	Name        string
	Description string
	ImgUrl      string

	BookDetail      *BookDetail
	EquipmentDetail *EquipmentDetail
}

type BookDetail struct {
	ISBNCode string
}

type EquipmentDetail struct {
	Count    int
	CountMax int
}

type ItemSearchQuery struct {
	Name string
	// TODO
}

type ItemRepository interface {
	GetByID(id int) (*Item, error)
	Search(query ItemSearchQuery) ([]*Item, error)
	Create(item *Item) error
	Update(item *Item) error
	Delete(id int) error
}
