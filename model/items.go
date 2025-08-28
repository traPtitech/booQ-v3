package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Item struct {
	GormModel
	Name                 string                 `gorm:"type:text;not null" json:"name"`
	Description          string                 `gorm:"type:text;" json:"description"`
	ImgURL               string                 `gorm:"type:text;" json:"imgUrl"`
	Book                 *Book                  `gorm:"foreignKey:item_id;references:id" json:"book,omitempty"`
	Equipment            *Equipment             `gorm:"foreignKey:item_id;references:id" json:"equipment,omitempty"`
	Comment              []Comment              `gorm:"foreignKey:item_id;references:id" json:"comment,omitempty"`
	Tag                  []Tag                  `gorm:"foreignKey:item_id;references:id" json:"tag,omitempty"`
	Ownership            []Ownership            `gorm:"foreignKey:item_id;references:id" json:"ownership,omitempty"`
	Like                 []Like                 `gorm:"foreignKey:item_id;references:id" json:"like,omitempty"`
	TransactionEquipment []TransactionEquipment `gorm:"foreignKey:item_id;references:id" json:"transactionEquipment,omitempty"`
}

func (Item) TableName() string {
	return "items"
}

type Book struct {
	GormModelWithoutID
	ItemID int    `gorm:"primaryKey" json:"itemId"`
	Code   string `gorm:"type:varchar(13);" json:"code"`
}

func (Book) TableName() string {
	return "books"
}

type Equipment struct {
	GormModelWithoutID
	ItemID   int `gorm:"primaryKey" json:"itemId"`
	Count    int `gorm:"type:int;not null" json:"count"`
	CountMax int `gorm:"type:int;not null" json:"countMax"`
}

func (Equipment) TableName() string {
	return "equipments"
}

type RequestPostItemsBody struct {
	Name        string   `json:"name"`
	IsTrapItem  bool     `json:"isTrapItem"`
	IsBook      bool     `json:"isBook"`
	Count       int      `json:"count"`
	Code        string   `json:"code"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	ImgURL      string   `json:"imgUrl"`
}

type GetItemsBody struct {
	UserID      string   `json:"userId"`
	Search      string   `json:"search"`
	Rental      string   `json:"rental"`
	Limit       int      `json:"limit"`
	Offset      int      `json:"offset"`
	Tags        []string `json:"tag"`
	TagsExclude []string `json:"tag-exclude"`
	SortBy      string   `json:"sortby"`
}

func dbPreloaded() *gorm.DB {
	return db.Preload("Book").Preload("Equipment").Preload("Comment").Preload("Tag").
		Preload("Ownership").Preload("Like").Preload("Ownership.Transaction").Preload("TransactionEquipment")
}

func GetItems(query GetItemsBody) ([]Item, error) {
	query.Limit = max(query.Limit, 20)

	model := db.Preload("Book").Preload("Equipment").Preload("Tag")
	model = model.Limit(query.Limit).Offset(query.Offset)

	if query.Search != "" {
		model = model.Where("name LIKE ?", "%"+query.Search+"%")
	}

	// TODO: userid, rental, tag, tag-exclude, sortby
	// sortby: 項目名 + 降順or昇順 の2つの情報が必要そう

	items := []Item{}
	if err := model.Find(&items).Error; err != nil {
		return []Item{}, err
	}
	return items, nil
}

func itemFromBody(itemBody RequestPostItemsBody) Item {
	item := Item{
		Name:        itemBody.Name,
		Description: itemBody.Description,
		ImgURL:      itemBody.ImgURL,
		Tag:         stringsToTags(itemBody.Tags),
	}
	if itemBody.IsBook {
		item.Book = &Book{Code: itemBody.Code}
	}
	if itemBody.IsTrapItem {
		item.Equipment = &Equipment{Count: itemBody.Count, CountMax: itemBody.Count}
	}
	return item
}

func CreateItems(itemBodies []RequestPostItemsBody, me string) ([]Item, error) {
	items := make([]Item, len(itemBodies))
	err := db.Transaction(func(tx *gorm.DB) error {
		for i, itemBody := range itemBodies {
			items[i] = itemFromBody(itemBody)

			if items[i].Equipment == nil {
				items[i].Ownership = []Ownership{{UserID: me, Rentalable: true}}
			}
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return []Item{}, err
	}

	return items, nil
}

func stringsToTags(tagStrs []string) []Tag {
	res := make([]Tag, len(tagStrs))
	for i, tagStr := range tagStrs {
		res[i] = Tag{Name: tagStr}
	}
	return res
}

func GetItem(itemID int) (Item, error) {
	res := Item{}
	if err := dbPreloaded().First(&res, itemID).Error; err != nil {
		return Item{}, err
	}
	return res, nil
}

// TODO: Booksの紐づけ, Ownershipsの除去/入力を適切に行う
// TODO: 備品がCount < CountMaxの際にPatchItemされたらどうする？
func PatchItem(itemID int, itemBody RequestPostItemsBody) (Item, error) {
	itemOld, err := GetItem(itemID)
	if err != nil {
		return Item{}, err
	}

	if itemOld.Book == nil && itemBody.IsBook || itemOld.Book != nil && !itemBody.IsBook {
		return Item{}, fmt.Errorf("それが本かどうかの情報を変えることはできません: %w", ErrUpdateNotAllowed)
	}
	if itemOld.Equipment == nil && itemBody.IsTrapItem || itemOld.Equipment != nil && !itemBody.IsTrapItem {
		return Item{}, fmt.Errorf("それが物品かどうかの情報を変えることはできません: %w", ErrUpdateNotAllowed)
	}

	item := itemFromBody(itemBody)
	item.Ownership = itemOld.Ownership
	if err := db.Model(&itemOld).Updates(item).Error; err != nil {
		return Item{}, err
	}

	return item, nil
}

func DeleteItem(itemID int) error {
	item, err := GetItem(itemID)
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if item.Ownership != nil {
			ownershipIDs := []int{}
			if err := tx.Model(&Ownership{}).Where("item_id = ?", item.ID).Pluck("id", &ownershipIDs).Error; err != nil {
				return err
			}
			if err := db.Where("ownership_id IN ?", ownershipIDs).Delete(&Transaction{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Select("Book", "Equipment", "Comment", "Tag", "Like", "Ownership",
			"TransactionEquipment").Delete(&item).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
