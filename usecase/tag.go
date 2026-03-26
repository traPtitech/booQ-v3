package usecase

import (
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
)

type TagUseCase interface {
	GetByItemID(itemID int) ([]*domain.Tag, error)
	ReplaceByItemID(itemID int, tags []string) error
}

type tagUseCase struct {
	tagRepo  domain.TagRepository
	itemRepo domain.ItemRepository
}

func NewTagUseCase(tagRepo domain.TagRepository, itemRepo domain.ItemRepository) TagUseCase {
	return &tagUseCase{
		tagRepo:  tagRepo,
		itemRepo: itemRepo,
	}
}

func (u *tagUseCase) GetByItemID(itemID int) ([]*domain.Tag, error) {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return nil, err
	}

	return u.tagRepo.GetByItemID(itemID)
}

func (u *tagUseCase) ReplaceByItemID(itemID int, tags []string) error {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return err
	}

	if err := u.tagRepo.ReplaceByItemID(itemID, tags); err != nil {
		return fmt.Errorf("failed to replace tags: %w", err)
	}

	return nil
}
