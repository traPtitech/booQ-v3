package model

type Ownership struct {
	GormModel
	ItemID      int           `gorm:"type:int;not null" json:"itemId"`
	UserID      string        `gorm:"type:varchar(32);not null" json:"userId"`
	Rentalable  bool          `gorm:"type:boolean;not null" json:"rentalable"`
	Memo        string        `gorm:"type:varchar(32)" json:"memo"`
	Transaction []Transaction `gorm:"foreignKey:ownership_id;references:id"`
}

func (Ownership) TableName() string {
	return "ownerships"
}
