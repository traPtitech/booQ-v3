package usecase

import (
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/traPtitech/booQ-v3/domain"
)

const (
	maxFileSize      = 3 * 1024 * 1024 // 3MB
	mimeTypeJPEG     = "image/jpeg"
	mimeTypePNG      = "image/png"
	fileExtensionJPG = ".jpg"
	fileExtensionPNG = ".png"
)

type FileUseCase interface {
	Upload(src io.Reader, contentType string, size int64) (*domain.File, error)
	GetFile(id int) (io.ReadCloser, *domain.File, error)
}

type fileUseCase struct {
	fileRepo    domain.FileRepository
	fileStorage domain.FileStorage
}

func NewFileUseCase(fileRepo domain.FileRepository, fileStorage domain.FileStorage) FileUseCase {
	return &fileUseCase{
		fileRepo:    fileRepo,
		fileStorage: fileStorage,
	}
}

func (u *fileUseCase) Upload(src io.Reader, contentType string, size int64) (*domain.File, error) {
	// バリデーション: ファイルサイズ
	if size > maxFileSize {
		return nil, domain.ErrFileTooLarge
	}

	// バリデーション: MIMEタイプ & 拡張子決定
	var ext string
	switch contentType {
	case mimeTypeJPEG:
		ext = fileExtensionJPG
	case mimeTypePNG:
		ext = fileExtensionPNG
	default:
		return nil, domain.ErrInvalidFileType
	}

	// ファイル名生成（UUID + 拡張子）
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// ストレージに保存
	if err := u.fileStorage.Save(filename, src); err != nil {
		return nil, err
	}

	// DBにメタデータ保存
	file, err := u.fileRepo.Create(&domain.File{
		Name:     filename,
		MimeType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (u *fileUseCase) GetFile(id int) (io.ReadCloser, *domain.File, error) {
	// DBからメタデータ取得
	file, err := u.fileRepo.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	// ストレージからファイル取得
	reader, err := u.fileStorage.Open(file.Name)
	if err != nil {
		return nil, nil, err
	}

	return reader, file, nil
}
