package usecase

import (
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
)

type OwnershipUseCase interface {
	GetByItemID(itemID int) ([]*domain.Ownership, error)
	CreateOwnership(ownership *domain.Ownership) (*domain.Ownership, error)
	UpdateOwnership(ownership *domain.Ownership, userID string) (*domain.Ownership, error)
	DeleteOwnership(id int, userID string) error
}

type ownershipUseCase struct {
	ownershipRepo domain.OwnershipRepository
}

func NewOwnershipUseCase(ownershipRepo domain.OwnershipRepository) OwnershipUseCase {
	return &ownershipUseCase{
		ownershipRepo: ownershipRepo,
	}
}

func (u *ownershipUseCase) GetByItemID(itemID int) ([]*domain.Ownership, error) {
	return u.ownershipRepo.GetByItemID(itemID)
}

func (u *ownershipUseCase) CreateOwnership(ownership *domain.Ownership) (*domain.Ownership, error) {
	return u.ownershipRepo.Create(ownership)
}

// TODO: 管理者も変えられるように
func (u *ownershipUseCase) UpdateOwnership(ownership *domain.Ownership, userID string) (*domain.Ownership, error) {
	o, err := u.ownershipRepo.GetByID(ownership.ID)
	if err != nil {
		return nil, err
	}

	if o.UserID != userID {
		return nil, fmt.Errorf("%w: you can only update your own ownership", ErrForbidden)
	}

	return u.ownershipRepo.Update(ownership)
}

func (u *ownershipUseCase) DeleteOwnership(id int, userID string) error {
	o, err := u.ownershipRepo.GetByID(id)
	if err != nil {
		return err
	}

	if o.UserID != userID {
		return fmt.Errorf("%w: you can only delete your own ownership", ErrForbidden)
	}

	return u.ownershipRepo.Delete(id)
}
