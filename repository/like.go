package repository

import (
	"errors"
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type like struct {
	GormModelWithoutID
	ItemID int    `gorm:"primarykey;type:int;not null"`
	UserID string `gorm:"primarykey;type:varchar(32);not null"`
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) domain.LikeRepository {
	return &likeRepository{db: db}
}

func (l *like) toDomain() *domain.Like {
	return &domain.Like{
		ItemID: l.ItemID,
		UserID: l.UserID,
	}
}

func (repo *likeRepository) GetByItemID(itemID int) ([]*domain.Like, error) {
	var models []like
	if err := repo.db.Where("item_id = ?", itemID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get likes by item id: %w", err)
	}

	res := make([]*domain.Like, len(models))
	for i, m := range models {
		res[i] = m.toDomain()
	}

	return res, nil
}

func (repo *likeRepository) Exists(itemID int, userID string) (bool, error) {
	var model like
	err := repo.db.Where("item_id = ? AND user_id = ?", itemID, userID).First(&model).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, fmt.Errorf("failed to check like existence: %w", err)
}

func (repo *likeRepository) Create(l *domain.Like) error {
	model := &like{ItemID: l.ItemID, UserID: l.UserID}
	if err := repo.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	return nil
}

func (repo *likeRepository) Delete(itemID int, userID string) error {
	res := repo.db.Where("item_id = ? AND user_id = ?", itemID, userID).Delete(&like{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete like: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
