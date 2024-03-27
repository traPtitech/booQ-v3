package model

type Like struct {
	GormModelWithoutID
	ItemID int    `gorm:"primaryKey;type:int;not null" json:"item_id"`
	UserID string `gorm:"primaryKey;type:varchar(32);not null" json:"user_id"`
}

func (Like) TableName() string {
	return "likes"
}
