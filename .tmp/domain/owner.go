package domain

type OwnerKind string

const (
	OwnerUser  OwnerKind = "user"
	OwnerGroup OwnerKind = "group"
)

// OwnerRef は参照値。User と Group で同じ ID でも別の所有者になる。
type OwnerRef struct {
	Kind OwnerKind
	ID   string
}

func (o OwnerRef) Validate() error {
	if o.ID == "" || (o.Kind != OwnerUser && o.Kind != OwnerGroup) {
		return ErrInvalidOwner
	}

	return nil
}

// User は別サービスのユーザーを表す参照モデル。ローカルのテーブルは持たない。
type User struct {
	TraQID string
}

// D08: 管理者情報は概念として保持するが、今回の認可には使わない。
type Group struct {
	ID       string
	AdminIDs []string
}
