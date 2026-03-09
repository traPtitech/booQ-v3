package repository

import (
	"errors"
	"time"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type equipmentTransaction struct {
	ID               int       `gorm:"primaryKey;autoIncrement"`
	ItemID           int       `gorm:"not null"`
	UserID           int       `gorm:"not null"`
	Purpose          string    `gorm:"type:text"`
	BorrowCount      int       `gorm:"not null"`
	DueDate          time.Time `gorm:"not null"`
	ReturnDate       *time.Time
	BorrowInClubRoom bool   `gorm:"not null"`
	Status           int    `gorm:"not null;default:1"` // 1=貸出中, 2=返却済
	ReturnMessage    string `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (equipmentTransaction) TableName() string {
	return "transaction_equipments"
}

type equipmentBorrowingRepository struct {
	db *gorm.DB
}

func NewEquipmentBorrowingRepository(db *gorm.DB) domain.EquipmentBorrowingRepository {
	return &equipmentBorrowingRepository{db: db}
}

func (t *equipmentTransaction) toDomain() *domain.EquipmentTransaction {
	return &domain.EquipmentTransaction{
		ID:               t.ID,
		ItemID:           t.ItemID,
		UserID:           t.UserID,
		Status:           t.Status,
		Purpose:          t.Purpose,
		DueDate:          t.DueDate,
		ReturnDate:       t.ReturnDate,
		BorrowCount:      t.BorrowCount,
		BorrowInClubRoom: t.BorrowInClubRoom,
		ReturnMessage:    t.ReturnMessage,
		CreatedAt:        t.CreatedAt,
		UpdatedAt:        t.UpdatedAt,
	}
}

func (r *equipmentBorrowingRepository) Borrow(itemID int, userID int, req domain.BorrowRequestEquipment) (*domain.EquipmentTransaction, error) {
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	var createdTx equipmentTransaction

	err = r.db.Transaction(func(tx *gorm.DB) error {

		createdTx = equipmentTransaction{
			ItemID:           itemID,
			UserID:           userID,
			Purpose:          req.Propose,
			BorrowCount:      req.Count,
			DueDate:          dueDate,
			BorrowInClubRoom: req.BorrowInClubRoom,
			Status:           1, // 1=貸出中
			ReturnDate:       nil,
		}

		if err := tx.Create(&createdTx).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTx.toDomain(), nil
}

func (r *equipmentBorrowingRepository) Return(itemID int, userID int, req domain.BorrowReturn) (*domain.EquipmentTransaction, error) {
	var targetTx equipmentTransaction

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("item_id = ? AND user_id = ? AND status = ?", itemID, userID, 1).
			Order("created_at DESC").
			First(&targetTx).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrBorrowingNotFound
			}
			return err
		}

		now := time.Now()
		targetTx.Status = 2 // 2=返却済み
		targetTx.ReturnDate = &now
		targetTx.ReturnMessage = req.Text

		if err := tx.Save(&targetTx).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return targetTx.toDomain(), nil
}
