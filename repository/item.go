package repository

import (
	"errors"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type item struct {
	GormModel
	Name        string `gorm:"type:text;not null"`
	Description string `gorm:"type:text"`
	ImgURL      string `gorm:"type:text"`
	// TODO
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) domain.ItemRepository {
	return &itemRepository{db: db}
}

func (i *item) toDomain() *domain.Item {
	return &domain.Item{
		ID:              i.ID,
		Name:            i.Name,
		Description:     i.Description,
		ImgUrl:          i.ImgURL,
		BookDetail:      nil,
		EquipmentDetail: nil,
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
	}
}

func (repo *itemRepository) GetByID(id int) (*domain.Item, error) {
	res := &item{}
	if err := repo.db.First(res, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrItemNotFound
		}
		return nil, err
	}

	return res.toDomain(), nil
}

func (repo *itemRepository) Search(query domain.ItemSearchQuery) ([]*domain.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *itemRepository) Create(item *domain.Item) error {
	//TODO implement me
	panic("implement me")
}

func (repo *itemRepository) Update(item *domain.Item) error {
	//TODO implement me
	panic("implement me")
}

func (repo *itemRepository) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
