package domain

import "time"

type Comment struct {
	ID        int
	ItemID    int
	UserID    string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CommentRepository interface {
	Create(comment *Comment) error
}
