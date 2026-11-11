package domain

import "time"

// D12: 本・備品の種別、ISBN、在庫を持たない共通の物品情報。
type Item struct {
	ID          int
	Name        string
	Description string
	ImgUrl      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type ItemSearchQuery struct {
	Name       string
	Owner      *OwnerRef
	BorrowerID string
	Limit      int
	Offset     int
	Tag        []string
	TagExclude []string
	SortBy     string
}

// 詳細表示用。Item 自体が全履歴を所有するという意味ではない。
type ItemDetail struct {
	Item       *Item
	Tags       []*Tag
	Likes      []*Like
	Ownerships []*OwnershipDetail
}

type OwnershipDetail struct {
	Ownership    *Ownership
	Requests     []*TransactionRequest
	Transactions []*Transaction
}

type ItemRepository interface {
	GetByID(id int) (*Item, error)
	GetDetailByID(id int) (*ItemDetail, error)
	Search(query ItemSearchQuery) ([]*ItemDetail, error)
	Create(item *Item) (*Item, error)
	CreateBatch(items []*Item) ([]*Item, error)
	Update(item *Item) (*Item, error)
	Delete(id int) error
}
