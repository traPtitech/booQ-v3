package model

import "fmt"

type Comment struct {
	GormModel
	ItemID  int    `gorm:"type:int;not null" json:"item_id"`
	UserID  string `gorm:"type:varchar(32);not null" json:"user_id"`
	Comment string `gorm:"type:text;not null" json:"comment"`
}

func (Comment) TableName() string {
	return "comments"
}

type CreateCommentPayload struct {
	ItemID  int    `json:"item_id"`
	UserID  string `json:"user_id"`
	Comment string `json:"comment"`
}

func CreateComment(p *CreateCommentPayload) (*Comment, error) {
	if p.ItemID == 0 {
		return nil, fmt.Errorf("ItemID is required")
	}
	if p.UserID == "" {
		return nil, fmt.Errorf("UserID is required")
	}
	if p.Comment == "" {
		return nil, fmt.Errorf("Comment is required")
	}
	c := &Comment{
		ItemID:  p.ItemID,
		UserID:  p.UserID,
		Comment: p.Comment,
	}
	if err := db.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}
