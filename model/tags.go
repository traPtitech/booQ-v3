package model

type Tag struct {
	GormModelWithoutID
	Name   string `gorm:"primaryKey;type:varchar(32);not null" json:"name"`
	ItemID int    `gorm:"primaryKey;type:int;not null" json:"itemId"`
}

func (Tag) TableName() string {
	return "tags"
}
