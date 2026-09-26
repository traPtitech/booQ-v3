package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/traPtitech/booQ-v3/internal/domain"
	"gorm.io/gorm"
)

type equipmentTransaction struct {
	GormModel
	ItemID        int        `gorm:"type:int;not null"`
	UserID        string     `gorm:"type:varchar(32);not null"`
	Purpose       string     `gorm:"type:text;not null"`
	Count         int        `gorm:"type:int;not null"`
	DueDate       time.Time  `gorm:"type:datetime;not null"`
	ReturnDate    *time.Time `gorm:"type:datetime"`
	Status        int        `gorm:"type:int;not null"` // 1=貸出中, 2=返却済
	ReturnMessage string     `gorm:"type:text;not null"`
}

func (equipmentTransaction) TableName() string {
	return "transactions_equipment"
}

func (t equipmentTransaction) toDomain() (*domain.EquipmentTransaction, error) {
	var status domain.EquipmentBorrowingStatus
	switch t.Status {
	case 1:
		status = domain.EquipmentBorrowingStatusBorrowed
	case 2:
		status = domain.EquipmentBorrowingStatusReturned
	default:
		return nil, fmt.Errorf("%w: unknown equipment status %d", domain.ErrInvalidTransactionStatus, t.Status)
	}
	result := &domain.EquipmentTransaction{
		ID:            t.ID,
		ItemID:        t.ItemID,
		UserID:        t.UserID,
		Status:        status,
		Purpose:       t.Purpose,
		Count:         t.Count,
		DueDate:       t.DueDate,
		ReturnMessage: t.ReturnMessage,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
	if t.ReturnDate != nil {
		result.ReturnDate = *t.ReturnDate
	}
	return result, nil
}

func toEquipmentTransactionModel(t *domain.EquipmentTransaction) (*equipmentTransaction, error) {
	var status int
	switch t.Status {
	case domain.EquipmentBorrowingStatusBorrowed:
		status = 1
	case domain.EquipmentBorrowingStatusReturned:
		status = 2
	default:
		return nil, fmt.Errorf("%w: unknown equipment status %q", domain.ErrInvalidTransactionStatus, t.Status)
	}
	model := &equipmentTransaction{
		GormModel:     GormModel{ID: t.ID, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt},
		ItemID:        t.ItemID,
		UserID:        t.UserID,
		Purpose:       t.Purpose,
		Count:         t.Count,
		DueDate:       t.DueDate,
		Status:        status,
		ReturnMessage: t.ReturnMessage,
	}
	if !t.ReturnDate.IsZero() {
		date := t.ReturnDate
		model.ReturnDate = &date
	}
	return model, nil
}

type equipmentTransactionRepository struct {
	db *gorm.DB
}

func NewEquipmentTransactionRepository(db *gorm.DB) domain.EquipmentTransactionRepository {
	return &equipmentTransactionRepository{db: db}
}

func (r *equipmentTransactionRepository) GetByID(id int) (*domain.EquipmentTransaction, error) {
	var model equipmentTransaction
	if err := r.db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return model.toDomain()
}

func (r *equipmentTransactionRepository) GetByUserID(userID string) ([]*domain.EquipmentTransaction, error) {
	return r.find(r.db.Where("user_id = ?", userID))
}

func (r *equipmentTransactionRepository) GetByItemID(itemID int) ([]*domain.EquipmentTransaction, error) {
	return r.find(r.db.Where("item_id = ?", itemID))
}

func (r *equipmentTransactionRepository) find(query *gorm.DB) ([]*domain.EquipmentTransaction, error) {
	var models []equipmentTransaction
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	results := make([]*domain.EquipmentTransaction, 0, len(models))
	for _, model := range models {
		transaction, err := model.toDomain()
		if err != nil {
			return nil, err
		}
		results = append(results, transaction)
	}
	return results, nil
}

func (r *equipmentTransactionRepository) Create(transaction *domain.EquipmentTransaction) (*domain.EquipmentTransaction, error) {
	model, err := toEquipmentTransactionModel(transaction)
	if err != nil {
		return nil, err
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	transaction.ID = model.ID
	transaction.CreatedAt = model.CreatedAt
	transaction.UpdatedAt = model.UpdatedAt
	return transaction, nil
}

func (r *equipmentTransactionRepository) Update(transaction *domain.EquipmentTransaction) (*domain.EquipmentTransaction, error) {
	model, err := toEquipmentTransactionModel(transaction)
	if err != nil {
		return nil, err
	}
	if err := r.db.Save(model).Error; err != nil {
		return nil, err
	}
	transaction.UpdatedAt = model.UpdatedAt
	return transaction, nil
}
