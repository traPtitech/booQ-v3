package usecase

import (
	"io"

	"github.com/traPtitech/booQ-v3/.tmp/domain"
)

type FileUseCase struct {
	Files   domain.FileRepository
	Storage domain.FileStorage
}

func (u FileUseCase) Upload(src io.Reader, contentType string, size int64) (*domain.File, error) {
	if size > 3*1024*1024 {
		return nil, domain.ErrFileTooLarge
	}

	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, domain.ErrInvalidFileType
	}

	// 以降は概念のみ（画像変換・UUID 生成・保存の実装は作らない）。
	// 1. src をデコードして向きを補正し、透過部分に白背景を付ける。
	// 2. JPEG に変換し、保存名を生成する。
	// 3. u.Storage.Save(filename, jpegReader) で本体を保存する。
	// 4. u.Files.Create(&domain.File{Name:filename, MimeType:"image/jpeg"})
	//    でメタデータを保存して返す。
	return nil, ErrConceptOnly
}

func (u FileUseCase) GetFile(id int) (io.ReadCloser, *domain.File, error) {
	file, err := u.Files.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	reader, err := u.Storage.Open(file.Name)
	if err != nil {
		return nil, nil, err
	}

	return reader, file, nil
}
