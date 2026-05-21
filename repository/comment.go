package repository

import (
	"time"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type comment struct {
	GormModel        // ID, CreatedAt, UpdatedAt
	ItemID    int    `gorm:"not null"`
	UserID    string `gorm:"type:varchar(32);not null"`
	Text      string `gorm:"column:comment;type:text;not null"`
	DeletedAt *time.Time
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(c *domain.Comment) (*domain.Comment, error) {
	var newComment *comment

	err := r.db.Transaction(func(tx *gorm.DB) error {
		newComment = &comment{
			ItemID: c.ItemID,
			UserID: c.UserID,
			Text:   c.Text,
		}

		if err := tx.Create(newComment).Error; err != nil {
			return err
		}

		return nil
	}).Error

	if err != nil {
		return nil, err
	}

	c.ID = newComment.ID
	c.CreatedAt = newComment.CreatedAt
	c.UpdatedAt = newComment.UpdatedAt
	c.DeletedAt = newComment.DeletedAt

	return c, nil
}
