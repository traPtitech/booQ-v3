package model

type Item struct {
	GormModel
	Name        string `gorm:"type:text;not null" json:"name"`
	Description string `gorm:"type:text;" json:"description"`
	ImgURL      string `gorm:"type:text;" json:"imgUrl"`
}

func (Item) TableName() string {
	return "items"
}

// TODO: 外部キー制約
type Book struct {
	GormModelWithoutID
	ItemID int    `gorm:"primary_key" json:"itemId"`
	Code   string `gorm:"type:varchar(13);" json:"code"`
}

func (Book) TableName() string {
	return "books"
}

type Equipment struct {
	GormModelWithoutID
	ItemID   int `gorm:"primary_key" json:"itemId"`
	Count    int `gorm:"type:int;not null" json:"count"`
	CountMax int `gorm:"type:int;not null" json:"countMax"`
}

func (Equipment) TableName() string {
	return "equipments"
}
