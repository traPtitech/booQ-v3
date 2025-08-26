package model

import (
	"gorm.io/gorm"
)

type Ownership struct {
	GormModel
	ItemID      int           `gorm:"type:int;not null" json:"itemId"`
	UserID      string        `gorm:"type:varchar(32);not null" json:"userId"`
	Rentalable  bool          `gorm:"type:boolean;not null" json:"rentalable"`
	Memo        string        `gorm:"type:varchar(32)" json:"memo"`
	Transaction []Transaction `gorm:"foreignKey:ownership_id;references:id"`
}

type OwnershipPayload struct {
	ItemID     int    `json:"itemId"`
	UserID     string `json:"userId"`
	Rentalable bool   `json:"rentalable"`
	Memo       string `json:"memo"`
}

func (Ownership) TableName() string {
	return "ownerships"
}

func CreateOwnership(ownership OwnershipPayload) error {
	o := Ownership{
		ItemID:     ownership.ItemID,
		UserID:     ownership.UserID,
		Rentalable: ownership.Rentalable,
		Memo:       ownership.Memo,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&o).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
