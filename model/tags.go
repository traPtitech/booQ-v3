package model

type Tag struct {
	GormModel
	Name   string `gorm:"type:text;not null" json:"name"`
	ItemID int    `gorm:"type:int;not null" json:"itemId"`
}

func (Tag) TableName() string {
	return "tags"
}
