package model

type Ownership struct {
	GormModel
	ItemID string `gorm:"type:int;not null" json:"name"`
	UserID string `gorm:"type:text;not null" json:"userId"`
}

func (Ownership) TableName() string {
	return "items"
}
