package model

import "time"

type Transaction struct {
	GormModel
	OwnershipID   int       `gorm:"type:int;not null" json:"ownershipId"`
	UserID        string    `gorm:"type:varchar(32);not null" json:"userId"`
	Status        int       `gorm:"type:int;not null" json:"status"`
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
	ItemID        int       `gorm:"type:int;not null" json:"itemId"`
	UserID        string    `gorm:"type:varchar(32);not null" json:"userId"`
	Status        int       `gorm:"type:int;not null" json:"status"`
	Purpose       string    `gorm:"type:text" json:"purpose"`
	ReturnMessage string    `gorm:"type:text" json:"returnMessage"`
	ReturnDue     time.Time `gorm:"type:datetime" json:"dueDate"`
	ReturnDate    time.Time `gorm:"type:datetime" json:"returnDate"`
}

func (TransactionEquipment) TableName() string {
	return "transactions_equipment"
}
