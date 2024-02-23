package model

type Comment struct {
	GormModel
	ItemID  int    `gorm:"type:int;not null" json:"item_id"`
	UserID  string `gorm:"type:varchar(32);not null" json:"user_id"`
	Comment string `gorm:"type:text;not null" json:"comment"`
}

func (Comment) TableName() string {
	return "comment"
}
