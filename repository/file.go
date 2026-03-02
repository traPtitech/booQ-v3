package repository

import (
	"errors"

	"github.com/traPtitech/booQ-v3/domain"
	"gorm.io/gorm"
)

type file struct {
	GormModel
	Name     string `gorm:"type:varchar(255);not null"` // UUID.拡張子
	MimeType string `gorm:"type:varchar(100);not null"` // image/jpeg, image/png
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) domain.FileRepository {
	return &fileRepository{db: db}
}

func (f *file) toDomain() *domain.File {
	return &domain.File{
		ID:        f.ID,
		Name:      f.Name,
		MimeType:  f.MimeType,
		CreatedAt: f.CreatedAt,
	}
}

func (repo *fileRepository) Create(domainFile *domain.File) error {
	f := &file{
		Name:     domainFile.Name,
		MimeType: domainFile.MimeType,
	}

	if err := repo.db.Create(f).Error; err != nil {
		return err
	}

	// IDをドメインモデルに反映
	domainFile.ID = f.ID
	domainFile.CreatedAt = f.CreatedAt

	return nil
}

func (repo *fileRepository) GetByID(id int) (*domain.File, error) {
	f := &file{}
	if err := repo.db.First(f, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return f.toDomain(), nil
}
