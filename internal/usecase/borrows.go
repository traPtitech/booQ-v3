package usecase

import (
	"fmt"
	"time"

	"github.com/traPtitech/booQ-v3/internal/domain"
)

type BorrowingUseCase interface {
	BorrowEquipment(itemID int, userID string, purpose string, count int, dueDate time.Time) (*domain.EquipmentTransaction, error)
	ReturnEquipment(itemID int, borrowingID int, userID string, message string) (*domain.EquipmentTransaction, error)
	PostRequest(userID string, ownershipID int, purpose string, dueDate time.Time, borrowInClubRoom bool) (*domain.Transaction, error)
	GetRequest(userID string, ownershipID int, borrowingID int) (*domain.Transaction, error)
	ReplyRequest(userID string, ownershipID int, borrowingID int, approve bool, message string) (*domain.Transaction, error)
	ReturnItem(userID string, ownershipID int, borrowingID int, message string) error
}

type borrowingUseCase struct {
	transactionRepo          domain.TransactionRepository
	ownershipRepo            domain.OwnershipRepository
	equipmentTransactionRepo domain.EquipmentTransactionRepository
	itemRepo                 domain.ItemRepository
}

func NewBorrowingUseCase(transactionRepo domain.TransactionRepository, ownershipRepo domain.OwnershipRepository, equipmentTransactionRepo domain.EquipmentTransactionRepository, itemRepo domain.ItemRepository) BorrowingUseCase {
	return &borrowingUseCase{
		transactionRepo:          transactionRepo,
		ownershipRepo:            ownershipRepo,
		equipmentTransactionRepo: equipmentTransactionRepo,
		itemRepo:                 itemRepo,
	}
}

func (b *borrowingUseCase) PostRequest(userID string, ownershipID int, purpose string, dueDate time.Time, borrowInClubRoom bool) (*domain.Transaction, error) {
	_, err := b.ownershipRepo.GetByID(ownershipID)
	if err != nil {
		return nil, fmt.Errorf("ownership with ID %d not found: %w", ownershipID, err)
	}

	if dueDate.Before(time.Now()) {
		return nil, ErrInvalidDueDate
	}

	t := domain.NewTransaction(userID, ownershipID, purpose, borrowInClubRoom, dueDate)
	return b.transactionRepo.Create(t)
}

func (b *borrowingUseCase) GetRequest(userID string, ownershipID int, borrowingID int) (*domain.Transaction, error) {
	t, err := b.transactionRepo.GetByID(borrowingID)
	if err != nil {
		return nil, fmt.Errorf("transaction with ID %d not found: %w", borrowingID, err)
	}

	if t.OwnershipID != ownershipID {
		return nil, fmt.Errorf("%w: transaction with ID %d does not belong to ownership with ID %d", domain.ErrNotFound, borrowingID, ownershipID)
	}

	if t.UserID != userID {
		return nil, fmt.Errorf("%w: transaction with ID %d does not belong to user with ID %s", ErrForbidden, borrowingID, userID)
	}

	return t, nil
}

func (b *borrowingUseCase) ReplyRequest(userID string, ownershipID int, borrowingID int, approve bool, message string) (*domain.Transaction, error) {
	t, err := b.transactionRepo.GetByID(borrowingID)
	if err != nil {
		return nil, fmt.Errorf("transaction with ID %d not found: %w", borrowingID, err)
	}

	if t.OwnershipID != ownershipID {
		return nil, fmt.Errorf("%w: transaction with ID %d does not belong to ownership with ID %d", domain.ErrNotFound, borrowingID, ownershipID)
	}

	o, err := b.ownershipRepo.GetByID(ownershipID)
	if err != nil {
		return nil, fmt.Errorf("ownership with ID %d not found: %w", ownershipID, err)
	}

	if o.UserID != userID {
		return nil, fmt.Errorf("%w: ownership with ID %d does not belong to user with ID %s", ErrForbidden, ownershipID, userID)
	}

	if approve {
		err = t.Approve(message)
	} else {
		err = t.Reject(message)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to reply request: %w", err)
	}

	return b.transactionRepo.Update(t)
}

func (b *borrowingUseCase) ReturnItem(userID string, ownershipID int, borrowingID int, message string) error {
	t, err := b.GetRequest(userID, ownershipID, borrowingID)
	if err != nil {
		return err
	}

	err = t.Return(message)
	if err != nil {
		return fmt.Errorf("failed to return item: %w", err)
	}

	_, err = b.transactionRepo.Update(t)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (b *borrowingUseCase) BorrowEquipment(itemID int, userID string, purpose string, count int, dueDate time.Time) (*domain.EquipmentTransaction, error) {
	item, err := b.itemRepo.GetByID(itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	if item.EquipmentDetail == nil {
		return nil, ErrItemNotEquipment
	}
	if !dueDate.After(time.Now()) {
		return nil, ErrInvalidDueDate
	}
	if count <= 0 {
		count = 1
	}

	return b.equipmentTransactionRepo.Create(domain.NewEquipmentTransaction(userID, itemID, purpose, count, dueDate))
}

func (b *borrowingUseCase) ReturnEquipment(itemID int, borrowingID int, userID string, message string) (*domain.EquipmentTransaction, error) {
	target, err := b.equipmentTransactionRepo.GetByID(borrowingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get equipment transaction: %w", err)
	}
	if target.ItemID != itemID {
		return nil, domain.ErrNotFound
	}
	if target.UserID != userID {
		return nil, ErrForbidden
	}
	if err := target.Return(message); err != nil {
		return nil, err
	}

	return b.equipmentTransactionRepo.Update(target)
}
