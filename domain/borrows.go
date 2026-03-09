package domain

import (
	"errors"
	"time"
)

var (
	ErrItemNotEquipment  = errors.New("item is not equipment")
	ErrNotEnoughStock    = errors.New("not enough stock")
	ErrBorrowingNotFound = errors.New("active borrowing not found")
)

type EquipmentTransaction struct {
	ID               int        `json:"id"`
	ItemID           int        `json:"itemId"`
	UserID           int        `json:"userId"`
	Status           int        `json:"status"` // 1:貸出中, 2:返却済
	Purpose          string     `json:"purpose"`
	DueDate          time.Time  `json:"due_date"`
	ReturnDate       *time.Time `json:"return_date"`
	BorrowCount      int        `json:"-"`
	BorrowInClubRoom bool       `json:"-"`
	ReturnMessage    string     `json:"-"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type BorrowRequestEquipment struct {
	Propose          string `json:"propose"`
	Count            int    `json:"count"`
	DueDate          string `json:"dueDate"`
	BorrowInClubRoom bool   `json:"borrowInClubRoom"`
}

type BorrowReturn struct {
	Text string `json:"text"`
}

type EquipmentBorrowingRepository interface {
	Borrow(itemID int, userID int, req BorrowRequestEquipment) (*EquipmentTransaction, error)
	Return(itemID int, userID int, req BorrowReturn) (*EquipmentTransaction, error)
}
