package model

import (
	"errors"

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

func (Ownership) TableName() string {
	return "ownerships"
}

func CreateOwnership(ownership Ownership) error {
	if ownership.ItemID == 0 {
		return errors.New("itemId is required")
	}
	if ownership.UserID == "" {
		return errors.New("userId is required")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ownership).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
