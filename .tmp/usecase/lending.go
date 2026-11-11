package usecase

import (
	"time"

	"github.com/traPtitech/booQ-v3/.tmp/domain"
)

type LendingUseCase struct {
	Ownerships   domain.OwnershipRepository
	Requests     domain.RequestRepository
	Transactions domain.TransactionRepository
}

// now はフローを説明するために引数で渡す。Clock や実行環境は作らない。
func (u LendingUseCase) Borrow(
	ownershipKey domain.OwnershipKey,
	borrowerID string,
	quantity int,
	purpose string,
	dueAt time.Time,
	borrowInClubRoom bool,
	now time.Time,
) (*domain.Transaction, error) {
	o, err := u.Ownerships.Get(ownershipKey)
	if err != nil {
		return nil, err
	}

	if err := o.RequirePolicy(domain.Immediate); err != nil {
		return nil, err
	}

	n, err := u.Transactions.SumBorrowedQuantity(o.Key)
	if err != nil {
		return nil, err
	}

	if err := o.CheckQuantity(quantity, n); err != nil {
		return nil, err
	}

	loan, err := domain.NewTransaction(
		ownershipKey,
		borrowerID,
		quantity,
		purpose,
		dueAt,
		borrowInClubRoom,
		nil,
		now,
	)
	if err != nil {
		return nil, err
	}

	// 残数確認〜作成は一つの業務操作。並行実行の制御は実装しない。
	return u.Transactions.Create(loan)
}

func (u LendingUseCase) RequestBorrow(
	ownershipKey domain.OwnershipKey,
	borrowerID string,
	quantity int,
	purpose string,
	dueAt time.Time,
	borrowInClubRoom bool,
	now time.Time,
) (*domain.TransactionRequest, error) {
	o, err := u.Ownerships.Get(ownershipKey)
	if err != nil {
		return nil, err
	}

	if err := o.RequirePolicy(domain.ApprovalRequired); err != nil {
		return nil, err
	}

	// D05: 非予約。現在の貸出中数量を引かず、所有総数だけ確認する。
	if err := o.CheckQuantity(quantity, 0); err != nil {
		return nil, err
	}

	request, err := domain.NewRequest(
		ownershipKey,
		borrowerID,
		quantity,
		purpose,
		dueAt,
		borrowInClubRoom,
		now,
	)
	if err != nil {
		return nil, err
	}

	return u.Requests.Create(request)
}

func (u LendingUseCase) ApproveRequest(actor string, id int, message string, now time.Time) (*domain.Transaction, error) {
	request, err := u.Requests.GetByID(id)
	if err != nil {
		return nil, err
	}

	if request.Status != domain.Pending {
		return nil, domain.ErrInvalidState
	}

	o, err := u.Ownerships.Get(request.OwnershipKey)
	if err != nil {
		return nil, err
	}

	if err := o.RequirePolicy(domain.ApprovalRequired); err != nil {
		return nil, err
	}

	n, err := u.Transactions.SumBorrowedQuantity(o.Key)
	if err != nil {
		return nil, err
	}

	if err := o.CheckQuantity(request.Quantity, n); err != nil {
		return nil, err
	}

	next := *request
	loan, err := next.Approve(actor, message, now)
	if err != nil {
		return nil, err
	}

	// 認可は全許可。申請承認と貸出作成を一体として扱うことだけ示す。
	// 残数確認も同じ業務操作の範囲。DB ロックや rollback は実装しない。
	return u.Requests.SaveApproval(&next, loan)
}

func (u LendingUseCase) RejectRequest(actor string, id int, message string, now time.Time) error {
	request, err := u.Requests.GetByID(id)
	if err != nil {
		return err
	}

	next := *request
	if err := next.Reject(actor, message, now); err != nil {
		return err
	}

	return u.Requests.Save(&next)
}

func (u LendingUseCase) CancelRequest(actor string, id int, now time.Time) error {
	request, err := u.Requests.GetByID(id)
	if err != nil {
		return err
	}

	next := *request
	if err := next.Cancel(actor, now); err != nil {
		return err
	}

	return u.Requests.Save(&next)
}

func (u LendingUseCase) Return(id int, message string, now time.Time) error {
	loan, err := u.Transactions.GetByID(id)
	if err != nil {
		return err
	}

	next := *loan

	// 現在の Policy と借り手本人かどうかは検査せず、全量返却する。
	if err := next.Return(message, now); err != nil {
		return err
	}

	return u.Transactions.Save(&next)
}

func (u LendingUseCase) GetRequest(id int) (*domain.TransactionRequest, error) {
	return u.Requests.GetByID(id)
}

func (u LendingUseCase) GetTransaction(id int) (*domain.Transaction, error) {
	return u.Transactions.GetByID(id)
}
