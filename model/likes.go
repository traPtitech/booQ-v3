package model

type Like struct {
	GormModel
	ItemId  int    `gorm:"type:int;not null" json:"item_id"`
	UserId  string `gorm:"type:text;not null" json:"user_id"`
}

func (Like) TableName() string {
	return "likes"
}