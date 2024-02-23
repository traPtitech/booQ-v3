package model

type Like struct {
	GormModel
	ItemID int    `gorm:"type:int;not null" json:"item_id"`
	UserID string `gorm:"type:text;not null" json:"user_id"`
}

func (Like) TableName() string {
	return "likes"
}
