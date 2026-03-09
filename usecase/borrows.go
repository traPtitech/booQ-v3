package usecase

import (
	"fmt"
	"time"

	"github.com/traPtitech/booQ-v3/domain"
)

type EquipmentBorrowingUseCase interface {
	BorrowEquipment(itemID int, userID int, req domain.BorrowRequestEquipment) (*domain.EquipmentTransaction, error)
	ReturnEquipment(itemID int, userID int, req domain.BorrowReturn) (*domain.EquipmentTransaction, error)
}

type equipmentBorrowingUseCase struct {
	repo domain.EquipmentBorrowingRepository
}

func NewEquipmentBorrowingUseCase(repo domain.EquipmentBorrowingRepository) EquipmentBorrowingUseCase {
	return &equipmentBorrowingUseCase{repo: repo}
}

func (u *equipmentBorrowingUseCase) BorrowEquipment(itemID int, userID int, req domain.BorrowRequestEquipment) (*domain.EquipmentTransaction, error) {
	if req.Count <= 0 {
		req.Count = 1
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due date format")
	}

	if dueDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, fmt.Errorf("due date cannot be in the past")
	}

	return u.repo.Borrow(itemID, userID, req)
}

func (u *equipmentBorrowingUseCase) ReturnEquipment(itemID int, userID int, req domain.BorrowReturn) (*domain.EquipmentTransaction, error) {
	if req.Text == "" {
		return nil, fmt.Errorf("return message is required")
	}

	return u.repo.Return(itemID, userID, req)
}
