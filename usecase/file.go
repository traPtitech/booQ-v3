package usecase

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

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

	// 先頭のバイトを読み込んで実際のMIMEタイプを検出
	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return nil, err
	}
	header = header[:n]

	// 実際のバイト列からMIMEタイプを検出
	detectedType := http.DetectContentType(header)

	// バリデーション: MIMEタイプ & 拡張子決定
	var ext string
	switch detectedType {
	case mimeTypeJPEG:
		ext = fileExtensionJPG
	case mimeTypePNG:
		ext = fileExtensionPNG
	default:
		return nil, domain.ErrInvalidFileType
	}

	// ファイル名生成（UUID + 拡張子）
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// 読み込んだヘッダーと残りのデータを結合
	reader := io.MultiReader(bytes.NewReader(header), src)

	// ストレージに保存
	if err := u.fileStorage.Save(filename, reader); err != nil {
		return nil, err
	}

	// DBにメタデータ保存
	file, err := u.fileRepo.Create(&domain.File{
		Name:     filename,
		MimeType: detectedType,
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
