package usecase

import (
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
)

type LikeUseCase interface {
	AddLike(itemID int, userID string) error
	RemoveLike(itemID int, userID string) error
	GetByItemID(itemID int) ([]*domain.Like, error)
	IsLiked(itemID int, userID string) (bool, error)
}

type likeUseCase struct {
	likeRepo domain.LikeRepository
	itemRepo domain.ItemRepository
}

func NewLikeUseCase(likeRepo domain.LikeRepository, itemRepo domain.ItemRepository) LikeUseCase {
	return &likeUseCase{
		likeRepo: likeRepo,
		itemRepo: itemRepo,
	}
}

func (u *likeUseCase) AddLike(itemID int, userID string) error {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return err
	}

	exists, err := u.likeRepo.Exists(itemID, userID)
	if err != nil {
		return fmt.Errorf("failed to check like: %w", err)
	}
	if exists {
		return ErrAlreadyLiked
	}

	if err := u.likeRepo.Create(&domain.Like{ItemID: itemID, UserID: userID}); err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	return nil
}

func (u *likeUseCase) RemoveLike(itemID int, userID string) error {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return err
	}

	exists, err := u.likeRepo.Exists(itemID, userID)
	if err != nil {
		return fmt.Errorf("failed to check like: %w", err)
	}
	if !exists {
		return ErrNotLiked
	}

	if err := u.likeRepo.Delete(itemID, userID); err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	return nil
}

func (u *likeUseCase) GetByItemID(itemID int) ([]*domain.Like, error) {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return nil, err
	}

	return u.likeRepo.GetByItemID(itemID)
}

func (u *likeUseCase) IsLiked(itemID int, userID string) (bool, error) {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return false, err
	}

	liked, err := u.likeRepo.Exists(itemID, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check like: %w", err)
	}

	return liked, nil
}
