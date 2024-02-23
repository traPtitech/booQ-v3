package model

type Ownership struct {
	GormModel
	ItemID      int           `gorm:"type:int;not null" json:"name"`
	UserID      string        `gorm:"type:varchar(32);not null" json:"userId"`
	Transaction []Transaction `gorm:"foreignKey:ownership_id;references:id"`
}

func (Ownership) TableName() string {
	return "ownerships"
}
