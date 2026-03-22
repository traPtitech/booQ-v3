package repository

import (
	"errors"
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type ownership struct {
	GormModel
	ItemID      int           `gorm:"type:int;not null"`
	UserID      string        `gorm:"type:text;not null"`
	Rentable    bool          `gorm:"type:boolean;not null"`
	Memo        string        `gorm:"type:text"`
	Transaction []transaction `gorm:"foreignKey:OwnershipID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ownershipRepository struct {
	db *gorm.DB
}

func NewOwnershipRepository(db *gorm.DB) domain.OwnershipRepository {
	return &ownershipRepository{db: db}
}

func (o *ownership) toDomain() *domain.Ownership {
	return &domain.Ownership{
		ID:       o.ID,
		ItemID:   o.ItemID,
		UserID:   o.UserID,
		Rentable: o.Rentable,
		Memo:     o.Memo,
	}
}

func toOwnershipModel(d *domain.Ownership) *ownership {
	return &ownership{
		GormModel: GormModel{ID: d.ID},
		ItemID:    d.ItemID,
		UserID:    d.UserID,
		Rentable:  d.Rentable,
		Memo:      d.Memo,
	}
}

func (repo *ownershipRepository) GetByID(id int) (*domain.Ownership, error) {
	res := &ownership{}

	if err := repo.db.First(res, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return res.toDomain(), nil
}

func (repo *ownershipRepository) GetByItemID(itemID int) ([]*domain.Ownership, error) {
	var models []ownership

	if err := repo.db.Where("item_id = ?", itemID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get ownerships by item id: %w", err)
	}

	res := make([]*domain.Ownership, len(models))
	for i, model := range models {
		res[i] = model.toDomain()
	}

	return res, nil
}

func (repo *ownershipRepository) GetByUserID(userID string) ([]*domain.Ownership, error) {
	var models []ownership

	if err := repo.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get ownerships by user id: %w", err)
	}

	res := make([]*domain.Ownership, len(models))
	for i, model := range models {
		res[i] = model.toDomain()
	}

	return res, nil
}

func (repo *ownershipRepository) Create(d *domain.Ownership) (*domain.Ownership, error) {
	model := toOwnershipModel(d)

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(model).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ownership: %w", err)
	}

	return model.toDomain(), nil
}

func (repo *ownershipRepository) Update(d *domain.Ownership) (*domain.Ownership, error) {
	model := toOwnershipModel(d)

	err := repo.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(model).Error
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update ownership: %w", err)
	}

	return model.toDomain(), nil
}

func (repo *ownershipRepository) Delete(id int) error {
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&ownership{}, id)
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
		return fmt.Errorf("failed to delete ownership: %w", err)
	}

	return nil
}
