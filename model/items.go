package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Item struct {
	GormModel
	Name                 string                 `gorm:"type:text;not null" json:"name"`
	Description          string                 `gorm:"type:text;" json:"description"`
	ImgURL               string                 `gorm:"type:text;" json:"imgUrl"`
	Book                 sql.Null[Book]         `gorm:"foreignKey:item_id;references:id" json:"book,omitempty"`
	Equipment            sql.Null[Equipment]    `gorm:"foreignKey:item_id;references:id" json:"equipment,omitempty"`
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

func GetItems(query GetItemsBody) ([]Item, error) {
	query.Limit = max(query.Limit, 20)

	model := db.Limit(query.Limit).Offset(query.Offset)
	// TODO: userid, rental, search, tag, tag-exclude, sortby

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
		item.Book = sql.Null[Book]{V: Book{Code: itemBody.Code}}
	} else {
		item.Book = sql.Null[Book]{Valid: false}
	}
	if itemBody.IsTrapItem {
		item.Equipment = sql.Null[Equipment]{V: Equipment{Count: itemBody.Count, CountMax: itemBody.Count}}
	} else {
		item.Equipment = sql.Null[Equipment]{Valid: false}
	}
	return item
}

func CreateItems(itemBodies []RequestPostItemsBody, me string) ([]Item, error) {
	items := make([]Item, len(itemBodies))
	err := db.Transaction(func(tx *gorm.DB) error {
		for i, itemBody := range itemBodies {
			items[i] = itemFromBody(itemBody)
			items[i].Ownership = []Ownership{{UserID: me, Rentalable: true}}
		}

		// TODO: アイテムの重複チェックができてるか要チェックしたいが
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
	if err := db.Preload("Book").Preload("Equipment").First(&res, itemID).Error; err != nil {
		return Item{}, err
	}
	return res, nil
}

func PatchItem(itemID int, itemBody RequestPostItemsBody) (Item, error) {
	itemOld, err := GetItem(itemID)
	if err != nil {
		return Item{}, err
	}

	item := itemFromBody(itemBody)
	if err := db.Model(&itemOld).Updates(item).Error; err != nil {
		return Item{}, err
	}

	return item, nil
}

func DeleteItem(itemID int) error {
	if err := db.Delete(&Item{}, itemID).Error; err != nil {
		return err
	}
	return nil
}
