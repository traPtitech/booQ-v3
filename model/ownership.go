package model

import (
	"errors"
	"fmt"

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

func CreateOwnership(ownership OwnershipPayload) (Ownership, error) {
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
		return Ownership{}, err
	}

	return o, nil
}

func UpdateOwnership(ownershipId int, ownership OwnershipPayload) (Ownership, error) {
	ownershipOld, err := GetOwnership(ownershipId)
	if err != nil {
		return Ownership{}, err
	}

	if ownershipOld.UserID != ownership.UserID {
		return Ownership{}, fmt.Errorf("編集する権限がありません: %w", ErrUnauthorized)
	}

	o := Ownership{
		ItemID:      ownership.ItemID,
		UserID:      ownership.UserID,
		Rentalable:  ownership.Rentalable,
		Memo:        ownership.Memo,
		Transaction: ownershipOld.Transaction,
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Ownership{}).Where("id = ?", ownershipId).Updates(map[string]interface{}{
			"rentalable": o.Rentalable,
			"memo":       o.Memo,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Ownership{}, err
	}

	return o, nil
}

func DeleteOwnership(ownershipId int, userId string) error {
	ownershipOld, err := GetOwnership(ownershipId)
	if err != nil {
		return err
	}

	if ownershipOld.UserID != userId {
		return fmt.Errorf("削除する権限がありません: %w", ErrUnauthorized)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("ownership_id = ?", ownershipId).Delete(&Transaction{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Ownership{}, ownershipId).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func GetOwnership(ownershipId int) (*Ownership, error) {
	var o Ownership
	if err := db.Preload("Transaction").First(&o, ownershipId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}
