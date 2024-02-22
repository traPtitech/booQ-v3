package model

import "time"

type Transaction struct {
	GormModel
	OwnershipId   int       `gorm:"type:int;not null" json:"ownershipId"`
	UserId        string    `gorm:"type:text;not null" json:"userId"`
	Statue        int       `gorm:"type:int;not null" json:"statue"`
	Purpose       string    `gorm:"type:text" json:"purpose"`
	Message       string    `gorm:"type:text" json:"message"`
	ReturnMessage string    `gorm:"type:text" json:"returnMessage"`
	ReturnDue     time.Time `gorm:"type:datetime" json:"dueDate"`
	CheckoutDue   time.Time `gorm:"type:datetime" json:"checkoutDate"`
	ReturnDate    time.Time `gorm:"type:datetime" json:"returnDate"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionEquipment struct {
	GormModel
	ItemId        int       `gorm:"type:int;not null" json:"itemId"`
	UserId        string    `gorm:"type:text;not null" json:"userId"`
	Statue        int       `gorm:"type:int;not null" json:"statue"`
	Purpose       string    `gorm:"type:text" json:"purpose"`
	ReturnMessage string    `gorm:"type:text" json:"returnMessage"`
	ReturnDue     time.Time `gorm:"type:datetime" json:"dueDate"`
	ReturnDate    time.Time `gorm:"type:datetime" json:"returnDate"`
}

func (TransactionEquipment) TableName() string {
	return "transactions_equipment"
}