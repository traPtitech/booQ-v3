package model

type Like struct {
	GormModelWithoutID
	ItemID int    `primaryKey;gorm:"type:int;not null" json:"item_id"`
	UserID string `primaryKey;gorm:"type:text;not null" json:"user_id"`
}

func (Like) TableName() string {
	return "likes"
}
