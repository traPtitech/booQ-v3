package model

type Ownership struct {
	GormModel
	ItemID      int           `gorm:"type:int;not null" json:"name"`
	UserID      string        `gorm:"type:varchar(32);not null" json:"userId"`
	Transaction []Transaction `gorm:"foreignKey:ownership_id;references:id"`
}

type PostOwnershipBody struct {
	UserID     string `json:"userId"`
	Rentalable bool   `json:"rentalable"`
	Memo       string `json:"memo"`
}

func (Ownership) TableName() string {
	return "ownerships"
}
