package repository

import (
	"errors"
	"fmt"

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

func toItemModel(d *domain.Item) *item {
	return &item{
		GormModel:   GormModel{ID: d.ID, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt},
		Name:        d.Name,
		Description: d.Description,
		ImgURL:      d.ImgUrl,
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

// TODO: search by tags, names, etc.
func (repo *itemRepository) Search(query domain.ItemSearchQuery) ([]*domain.Item, error) {
	var items []item
	dbQuery := repo.db.Model(&item{})

	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Limit > 0 {
		dbQuery = dbQuery.Limit(query.Limit)
	}
	if query.Offset > 0 {
		dbQuery = dbQuery.Offset(query.Offset)
	}

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to search items: %w", err)
	}

	domainItems := make([]*domain.Item, len(items))
	for i, item := range items {
		domainItems[i] = item.toDomain()
	}

	return domainItems, nil
}

func (repo *itemRepository) Create(item *domain.Item) (*domain.Item, error) {
	model := toItemModel(item)

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(model).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	item.ID = model.ID
	item.CreatedAt = model.CreatedAt
	item.UpdatedAt = model.UpdatedAt

	return item, nil
}

func (repo *itemRepository) CreateBatch(items []*domain.Item) ([]*domain.Item, error) {
	if len(items) == 0 {
		return items, nil
	}

	models := make([]*item, 0, len(items))
	for _, item := range items {
		models = append(models, toItemModel(item))
	}

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(models).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create items in batch: %w", err)
	}

	for i, model := range models {
		items[i].ID = model.ID
		items[i].CreatedAt = model.CreatedAt
		items[i].UpdatedAt = model.UpdatedAt
	}

	return items, nil
}

func (repo *itemRepository) Update(item *domain.Item) (*domain.Item, error) {
	model := toItemModel(item)

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(model).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	item.ID = model.ID
	item.CreatedAt = model.CreatedAt
	item.UpdatedAt = model.UpdatedAt

	return item, nil
}

func (repo *itemRepository) Delete(id int) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&item{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrItemNotFound
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
