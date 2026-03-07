package usecase

import (
	"errors"
	"fmt"

	"github.com/traPtitech/booQ-v3/domain"
)

type ItemUseCase interface {
	GetItemByID(id int) (*domain.Item, error)
	SearchItems(query domain.ItemSearchQuery) ([]*domain.Item, error)
	CreateItem(item *domain.Item) (*domain.Item, error)
	CreateItems(items []*domain.Item) ([]*domain.Item, error)
	UpdateItem(item *domain.Item) (*domain.Item, error)
	DeleteItem(id int) error
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

func (u *itemUseCase) SearchItems(query domain.ItemSearchQuery) ([]*domain.Item, error) {
	if query.Offset > 0 && query.Limit <= 0 {
		return nil, fmt.Errorf("%w: offset is set but limit is not set", ErrInvalidSearchQuery)
	}

	return u.itemRepo.Search(query)
}

func (u *itemUseCase) CreateItem(item *domain.Item) (*domain.Item, error) {
	return u.itemRepo.Create(item)
}

func (u *itemUseCase) CreateItems(items []*domain.Item) ([]*domain.Item, error) {
	created, err := u.itemRepo.CreateBatch(items)
	if err != nil {
		return nil, fmt.Errorf("failed to create items: %w", err)
	}

	return created, nil
}

// TODO: updateの認可
func (u *itemUseCase) UpdateItem(item *domain.Item) (*domain.Item, error) {
	itemOld, err := u.itemRepo.GetByID(item.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if itemOld.BookDetail == nil && item.BookDetail != nil || itemOld.BookDetail != nil && item.BookDetail == nil {
		return nil, fmt.Errorf("%w: cannot change whether item is book or not", ErrUpdateNotAllowed)
	}
	if itemOld.EquipmentDetail == nil && item.EquipmentDetail != nil || itemOld.EquipmentDetail != nil && item.EquipmentDetail == nil {
		return nil, fmt.Errorf("%w: cannot change whether item is equipment or not", ErrUpdateNotAllowed)
	}

	return u.itemRepo.Update(item)
}

// TODO: こっちも認可
func (u *itemUseCase) DeleteItem(id int) error {
	return u.itemRepo.Delete(id)
}
