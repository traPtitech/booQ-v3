package model

type Item struct {
	GormModel
	Name                 string                 `gorm:"type:text;not null" json:"name"`
	Description          string                 `gorm:"type:text;" json:"description"`
	ImgURL               string                 `gorm:"type:text;" json:"imgUrl"`
	Book                 Book                   `gorm:"foreignKey:item_id;references:id"`
	Equipment            Equipment              `gorm:"foreignKey:item_id;references:id"`
	Comment              []Comment              `gorm:"foreignKey:item_id;references:id"`
	Tag                  []Tag                  `gorm:"foreignKey:item_id;references:id"`
	Ownership            []Ownership            `gorm:"foreignKey:item_id;references:id"`
	Like                 []Like                 `gorm:"foreignKey:item_id;references:id"`
	TransactionEquipment []TransactionEquipment `gorm:"foreignKey:item_id;references:id"`
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

func GetItemByID(id int) (Item, error) {
	res := Item{}
	if err := db.First(&res, id).Error; err != nil {
		return Item{}, err
	}

	return res, nil
}
