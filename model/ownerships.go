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

type PostOwnershipBody struct {
	Rentalable bool   `json:"rentalable"`
	Memo       string `json:"memo"`
}

func (Ownership) TableName() string {
	return "ownerships"
}

func RegisterOwnership(ownership Ownership) (Ownership, error) {
	item, err := GetItem(ownership.ItemID)
	if err != nil {
		return Ownership{}, err
	}
	if item.Equipment != nil {
		return Ownership{}, errors.New("備品に対して所有者を設定することはできません。物品の情報を変更することで所有物の個数を変更してください")
	}
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
	// ownershipが存在するか確認
	if err := db.First(&Ownership{}, ownership.ID).Error; err != nil {
		return err
	}

	if err := db.Save(&ownership).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOwnership(id int) error {
	// ownershipが存在するか確認
	if err := db.First(&Ownership{}, id).Error; err != nil {
		return err
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 関連するトランザクションを削除
		if err := tx.Delete(&Transaction{}, "ownership_id = ?", id).Error; err != nil {
			return err
		}

		// 所有権を削除
		if err := tx.Delete(&Ownership{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}
