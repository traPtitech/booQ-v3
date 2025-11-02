package repository

import "github.com/traPtitech/booQ-v3/domain"

type item struct {
	GormModel
	Name        string `gorm:"type:text;not null"`
	Description string `gorm:"type:text"`
	ImgURL      string `gorm:"type:text"`
	// TODO
}

func (i *item) toDomain() *domain.Item {
	return &domain.Item{
		ID:              i.ID,
		Name:            i.Name,
		Description:     i.Description,
		ImgUrl:          i.ImgURL,
		BookDetail:      nil,
		EquipmentDetail: nil,
		CreatedAt:       i.CreatedAt,
		UpdatedAt:       i.UpdatedAt,
	}
}

func (d *DB) GetByID(id int) (*domain.Item, error) {
	res := &item{}
	if err := d.db.First(res, id).Error; err != nil {
		return nil, err
	}

	return res.toDomain(), nil
}

func (d *DB) Search(query domain.ItemSearchQuery) ([]*domain.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DB) Create(item *domain.Item) error {
	//TODO implement me
	panic("implement me")
}

func (d *DB) Update(item *domain.Item) error {
	//TODO implement me
	panic("implement me")
}

func (d *DB) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
