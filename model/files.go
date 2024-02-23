package model

import (
	"fmt"

	"io"

	"github.com/traPtitech/booQ-v3/storage"
	"gorm.io/gorm"
)

// File アップロードファイルの構造体
type File struct {
	GormModel
	UploadUserID uint `gorm:"type:int;not null"`
}

// TableName dbのテーブル名を指定する
func (File) TableName() string {
	return "files"
}

// CreateFile Fileを作成する
func CreateFile(uploadUserID uint, src io.Reader, ext string) (File, error) {
	f := File{UploadUserID: uploadUserID}
	// トランザクション開始
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&f).Error; err != nil {
			return err
		}

		if err := storage.Save(fmt.Sprintf("%d.%s", f.ID, ext), src); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return File{}, err
	}

	return f, nil
}
