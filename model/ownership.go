package model

type Ownership struct {
	GormModel
	ItemID      int           `gorm:"type:int;not null" json:"itemId"`
	UserID      string        `gorm:"type:varchar(32);not null" json:"userId"`
	Rentalable  bool          `gorm:"type:boolean;not null" json:"rentalable"`
	Memo        string        `gorm:"type:varchar(32)" json:"memo"`
	Transaction []Transaction `gorm:"foreignKey:ownership_id;references:id"`
}

type PostOwnershipBody struct {
	Rentalable bool   `json:"rentalable"`
	Memo       string `json:"memo"`
}

func (Ownership) TableName() string {
	return "ownerships"
}

func RegisterOwnership(ownership Ownership) (Ownership, error) {
	if err := db.Create(&ownership).Error; err != nil {
		return Ownership{}, nil
	}
	return ownership, nil
}

func GetOwnershipByID(id int) (Ownership, error) {
	res := Ownership{}
	if err := db.First(&res, id).Error; err != nil {
		return Ownership{}, err
	}
	return res, nil
}

func PatchOwnership(ownership Ownership) error {
	if err := db.Save(&ownership).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOwnership(id int) error {
	if err := db.Delete(&Ownership{}, id).Error; err != nil {
		return err
	}
	return nil
}
