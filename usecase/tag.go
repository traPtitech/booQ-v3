package usecase

import (
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
)

type TagUseCase interface {
	GetByItemID(itemID int) ([]*domain.Tag, error)
	GetByItemIDs(itemIDs []int) (map[int][]*domain.Tag, error)
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

func (u *tagUseCase) GetByItemIDs(itemIDs []int) (map[int][]*domain.Tag, error) {
	return u.tagRepo.GetByItemIDs(itemIDs)
}

func (u *tagUseCase) ReplaceByItemID(itemID int, tags []string) error {
	if _, err := u.itemRepo.GetByID(itemID); err != nil {
		return err
	}

	uniqueTags := make(map[string]struct{})
	for _, tag := range tags {
		uniqueTags[tag] = struct{}{}
	}

	uniqueTagList := make([]string, 0, len(uniqueTags))
	for tag := range uniqueTags {
		uniqueTagList = append(uniqueTagList, tag)
	}

	if err := u.tagRepo.ReplaceByItemID(itemID, uniqueTagList); err != nil {
		return fmt.Errorf("failed to replace tags: %w", err)
	}

	return nil
}
