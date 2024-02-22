package model

type Comment struct {
	GormModel
	ItemId  int    `gorm:"type:int;not null" json:"item_id"`
	UserId  string `gorm:"type:text;not null" json:"user_id"`
	Comment string `gorm:"type:text;not null" json:"comment"`
}

func (Comment) TableName() string {
	return "comment"
}