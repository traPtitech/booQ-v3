package repository

import (
	"errors"
	"time"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type transaction struct {
	GormModel
	ItemID           int        `gorm:"type:int;not null"`
	UserID           string     `gorm:"type:varchar(32);not null"`
	OwnershipID      int        `gorm:"type:int;not null"`
	Status           string     `gorm:"type:varchar(32);not null"`
	Purpose          string     `gorm:"type:text;not null"`
	BorrowInClubRoom bool       `gorm:"type:boolean;not null"`
	DueDate          time.Time  `gorm:"type:timestamp;not null"`
	Message          string     `gorm:"type:text;not null"`
	CheckoutDate     *time.Time `gorm:"type:timestamp"`
	ReturnMessage    string     `gorm:"type:text;not null"`
	ReturnDate       *time.Time `gorm:"type:timestamp"`
}

func (t transaction) toDomain() *domain.Transaction {
	return &domain.Transaction{
		ID:               t.ID,
		ItemID:           t.ItemID,
		UserID:           t.UserID,
		OwnershipID:      t.OwnershipID,
		Status:           domain.BorrowingStatus(t.Status),
		Purpose:          t.Purpose,
		BorrowInClubRoom: t.BorrowInClubRoom,
		DueDate:          t.DueDate,
		Message:          t.Message,
		CheckoutDate:     t.CheckoutDate,
		ReturnMessage:    t.ReturnMessage,
		ReturnDate:       t.ReturnDate,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func toTransactionModel(t *domain.Transaction) *transaction {
	return &transaction{
		GormModel:        GormModel{ID: t.ID, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt},
		ItemID:           t.ItemID,
		UserID:           t.UserID,
		OwnershipID:      t.OwnershipID,
		Status:           string(t.Status),
		Purpose:          t.Purpose,
		BorrowInClubRoom: t.BorrowInClubRoom,
		DueDate:          t.DueDate,
		Message:          t.Message,
		CheckoutDate:     t.CheckoutDate,
		ReturnMessage:    t.ReturnMessage,
		ReturnDate:       t.ReturnDate,
	}
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (t transactionRepository) GetByID(id int) (*domain.Transaction, error) {
	var transaction transaction
	if err := t.db.First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return transaction.toDomain(), nil
}

func (t transactionRepository) GetByUserID(userID string) ([]*domain.Transaction, error) {
	var transactions []transaction
	if err := t.db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	domainTransactions := make([]*domain.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		domainTransactions = append(domainTransactions, transaction.toDomain())
	}
	return domainTransactions, nil
}

func (t transactionRepository) GetByOwnershipID(ownershipID int) ([]*domain.Transaction, error) {
	var transactions []transaction
	if err := t.db.Where("ownership_id = ?", ownershipID).Find(&transactions).Error; err != nil {
		return nil, err
	}

	domainTransactions := make([]*domain.Transaction, 0, len(transactions))
	for _, transaction := range transactions {
		domainTransactions = append(domainTransactions, transaction.toDomain())
	}
	return domainTransactions, nil
}

func (t transactionRepository) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	model := toTransactionModel(transaction)
	err := t.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(model).Error
	})
	if err != nil {
		return nil, err
	}

	transaction.ID = model.ID
	transaction.CreatedAt = model.CreatedAt
	transaction.UpdatedAt = model.UpdatedAt
	return transaction, nil
}

func (t transactionRepository) Update(transaction *domain.Transaction) (*domain.Transaction, error) {
	model := toTransactionModel(transaction)
	err := t.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(model).Error
	})
	if err != nil {
		return nil, err
	}

	transaction.ID = model.ID
	transaction.CreatedAt = model.CreatedAt
	transaction.UpdatedAt = model.UpdatedAt
	return transaction, nil
}
