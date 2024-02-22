package model

type Ownership struct {
	GormModel
	ItemId      string `gorm:"type:int;not null" json:"name"`
	UserId      string `gorm:"type:text;not null" json:"userId"`
}

func (Ownership) TableName() string {
	return "items"
}