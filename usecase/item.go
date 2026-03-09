package usecase

import "github.com/traPtitech/booQ-v3/domain"

type ItemUseCase interface {
	GetItemByID(id int) (*domain.Item, error)
	// TODO: other methods
}

type itemUseCase struct {
	itemRepo domain.ItemRepository
}

func NewItemUseCase(itemRepo domain.ItemRepository) ItemUseCase {
	return &itemUseCase{
		itemRepo: itemRepo,
	}
}

func (u *itemUseCase) GetItemByID(id int) (*domain.Item, error) {
	return u.itemRepo.GetByID(id)
}
