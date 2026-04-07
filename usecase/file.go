package usecase

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/traPtitech/booQ-v3/domain"
)

const maxFileSize = 3 * 1024 * 1024 // 3MB

// アップロード可能なMIMEタイプ
var uploadableMimes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

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

	// バリデーション: MIMEタイプ
	if !uploadableMimes[contentType] {
		return nil, domain.ErrInvalidFileType
	}

	// 画像をデコード（AutoOrientationでEXIF情報による向き補正）
	orig, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return nil, domain.ErrInvalidFileType
	}

	// 新しいRGBA画像を作成し、白背景を塗る（PNG透過対応）
	newImg := image.NewRGBA(orig.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), orig, orig.Bounds().Min, draw.Over)

	// JPEGにエンコード
	buf := &bytes.Buffer{}
	if err := imaging.Encode(buf, newImg, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
		return nil, err
	}

	// ファイル名生成（UUID + .jpg）
	filename := fmt.Sprintf("%s.jpg", uuid.New().String())

	// ストレージに保存
	if err := u.fileStorage.Save(filename, buf); err != nil {
		return nil, err
	}

	// DBにメタデータ保存
	file, err := u.fileRepo.Create(&domain.File{
		Name:     filename,
		MimeType: "image/jpeg",
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
