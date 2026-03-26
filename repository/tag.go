package repository

import (
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type tag struct {
	GormModelWithoutID
	ItemID int    `gorm:"primarykey;type:int;not null;index"`
	Name   string `gorm:"primarykey;type:varchar(64);not null"`
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) domain.TagRepository {
	return &tagRepository{db: db}
}

func (t *tag) toDomain() *domain.Tag {
	return &domain.Tag{
		ItemID: t.ItemID,
		Name:   t.Name,
	}
}

func (repo *tagRepository) GetByItemID(itemID int) ([]*domain.Tag, error) {
	var models []tag
	if err := repo.db.Where("item_id = ?", itemID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get tags by item id: %w", err)
	}

	res := make([]*domain.Tag, len(models))
	for i, m := range models {
		res[i] = m.toDomain()
	}

	return res, nil
}

func (repo *tagRepository) ReplaceByItemID(itemID int, tags []string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("item_id = ?", itemID).Delete(&tag{}).Error; err != nil {
			return fmt.Errorf("failed to delete tags: %w", err)
		}

		if len(tags) == 0 {
			return nil
		}

		models := make([]*tag, 0, len(tags))
		for _, name := range tags {
			models = append(models, &tag{ItemID: itemID, Name: name})
		}

		if err := tx.Create(models).Error; err != nil {
			return fmt.Errorf("failed to create tags: %w", err)
		}

		return nil
	})
}
