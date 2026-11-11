package usecase

import "github.com/traPtitech/booQ-v3/.tmp/domain"

type OwnershipUseCase struct {
	Items        domain.ItemRepository
	Ownerships   domain.OwnershipRepository
	Requests     domain.RequestRepository
	Transactions domain.TransactionRepository
}

func (u OwnershipUseCase) ListOwnerships(itemID int) ([]*domain.Ownership, error) {
	return u.Ownerships.ListByItemID(itemID)
}

func (u OwnershipUseCase) CreateOwnership(o domain.Ownership) (*domain.Ownership, error) {
	if _, err := u.Items.GetByID(o.Key.ItemID); err != nil {
		return nil, err
	}

	if err := o.Validate(); err != nil {
		return nil, err
	}

	// User / Group の存在解決は外部の関心事。認可は全許可。
	return u.Ownerships.Create(&o)
}

func (u OwnershipUseCase) GetAvailability(key domain.OwnershipKey) (domain.Availability, error) {
	o, err := u.Ownerships.Get(key)
	if err != nil {
		return domain.Availability{}, err
	}

	n, err := u.Transactions.SumBorrowedQuantity(key)
	if err != nil {
		return domain.Availability{}, err
	}

	return o.Availability(n), nil
}

type UpdateOwnershipCommand struct {
	Owner    domain.OwnerRef
	Quantity int
	Policy   domain.LendingPolicy
	Memo     string
}

func (u OwnershipUseCase) UpdateOwnership(key domain.OwnershipKey, cmd UpdateOwnershipCommand) error {
	o, err := u.Ownerships.Get(key)
	if err != nil {
		return err
	}

	if cmd.Owner != key.Owner() {
		// TODO(design D01/D09/D11): 所有者変更でキーも変わる。
		// 過去の参照維持と移転先に同じキーがある場合の統合方針は未決。
		return ErrDesignPending
	}

	pending, err := u.Requests.HasPending(key)
	if err != nil {
		return err
	}

	if pending && cmd.Policy != o.Policy {
		// TODO(design D11): 既存申請の維持・取消を決めてから変更を記述する。
		return ErrDesignPending
	}

	borrowed, err := u.Transactions.SumBorrowedQuantity(key)
	if err != nil {
		return err
	}

	next := *o
	if err := next.ChangeSettings(cmd.Quantity, cmd.Policy, cmd.Memo, borrowed); err != nil {
		return err
	}

	return u.Ownerships.Save(&next)
}

func (u OwnershipUseCase) EnsureUnused(key domain.OwnershipKey) error {
	borrowed, err := u.Transactions.SumBorrowedQuantity(key)
	if err != nil {
		return err
	}

	pending, err := u.Requests.HasPending(key)
	if err != nil {
		return err
	}

	if borrowed > 0 || pending {
		return domain.ErrOwnershipInUse
	}

	return nil
}

func (u OwnershipUseCase) DeleteOwnership(key domain.OwnershipKey) error {
	if _, err := u.Ownerships.Get(key); err != nil {
		return err
	}

	if err := u.EnsureUnused(key); err != nil {
		return err
	}

	return u.Ownerships.Delete(key)
}
