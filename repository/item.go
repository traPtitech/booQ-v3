package repository

import (
	"errors"
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type item struct {
	GormModel
	Name        string     `gorm:"type:text;not null"`
	Description string     `gorm:"type:text"`
	ImgURL      string     `gorm:"type:text"`
	Book        *book      `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Equipment   *equipment `gorm:"foreignKey:ItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type book struct {
	GormModelWithoutID
	ItemID   int    `gorm:"primaryKey"`
	ISBNCode string `gorm:"type:varchar(13);"`
}

type equipment struct {
	GormModelWithoutID
	ItemID   int `gorm:"primaryKey"`
	Count    int `gorm:"type:int;not null"`
	CountMax int `gorm:"type:int;not null"`
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) domain.ItemRepository {
	return &itemRepository{db: db}
}

func (i *item) toDomain() *domain.Item {
	item := &domain.Item{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		ImgUrl:      i.ImgURL,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
	if i.Book != nil {
		item.BookDetail = &domain.BookDetail{
			ISBNCode: i.Book.ISBNCode,
		}
	}
	if i.Equipment != nil {
		item.EquipmentDetail = &domain.EquipmentDetail{
			Count:    i.Equipment.Count,
			CountMax: i.Equipment.CountMax,
		}
	}
	return item
}

func toItemModel(d *domain.Item) *item {
	item := &item{
		GormModel:   GormModel{ID: d.ID, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt},
		Name:        d.Name,
		Description: d.Description,
		ImgURL:      d.ImgUrl,
	}
	if d.BookDetail != nil {
		item.Book = &book{
			ItemID:   d.ID,
			ISBNCode: d.BookDetail.ISBNCode,
		}
	}
	if d.EquipmentDetail != nil {
		item.Equipment = &equipment{
			ItemID:   d.ID,
			Count:    d.EquipmentDetail.Count,
			CountMax: d.EquipmentDetail.CountMax,
		}
	}
	return item
}

func (repo *itemRepository) GetByID(id int) (*domain.Item, error) {
	res := &item{}

	model := repo.db.Preload("Book").Preload("Equipment").Model(&item{})
	if err := model.First(res, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return res.toDomain(), nil
}

// TODO: search by tags, names, etc.
func (repo *itemRepository) Search(query domain.ItemSearchQuery) ([]*domain.Item, error) {
	var items []item
	dbQuery := repo.db.Preload("Book").Preload("Equipment").Model(&item{})

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
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(model).Error
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
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
