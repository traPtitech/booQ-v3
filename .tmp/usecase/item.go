package usecase

import "github.com/traPtitech/booQ-v3/.tmp/domain"

type ItemUseCase struct {
	Items      domain.ItemRepository
	Ownerships OwnershipUseCase
}

func (u ItemUseCase) GetItemByID(id int) (*domain.Item, error) {
	return u.Items.GetByID(id)
}

func (u ItemUseCase) GetItemDetailByID(id int) (*domain.ItemDetail, error) {
	return u.Items.GetDetailByID(id)
}

func (u ItemUseCase) SearchItems(q domain.ItemSearchQuery) ([]*domain.ItemDetail, error) {
	if q.Offset > 0 && q.Limit <= 0 {
		return nil, ErrInvalidSearchQuery
	}

	return u.Items.Search(q)
}

func (u ItemUseCase) CreateItem(item *domain.Item) (*domain.Item, error) {
	return u.Items.Create(item)
}

func (u ItemUseCase) CreateItems(items []*domain.Item) ([]*domain.Item, error) {
	return u.Items.CreateBatch(items)
}

func (u ItemUseCase) UpdateItem(item *domain.Item) (*domain.Item, error) {
	if _, err := u.Items.GetByID(item.ID); err != nil {
		return nil, err
	}

	// 所有者や数量は UpdateOwnership の担当。物品種別の変更制限も不要。
	return u.Items.Update(item)
}

func (u ItemUseCase) DeleteItem(id int) error {
	if _, err := u.Items.GetByID(id); err != nil {
		return err
	}

	owners, err := u.Ownerships.ListOwnerships(id)
	if err != nil {
		return err
	}

	for _, o := range owners {
		if err := u.Ownerships.EnsureUnused(o.Key); err != nil {
			return err
		}
	}

	// 履歴の保持方法・関連する所有の廃止処理は repository 側の説明に留める。
	return u.Items.Delete(id)
}
