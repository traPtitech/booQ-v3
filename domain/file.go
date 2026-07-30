package domain

import (
	"io"
	"time"
)

type File struct {
	ID        int
	Name      string // 保存時のファイル名（UUID等）
	MimeType  string // image/jpeg, image/png
	CreatedAt time.Time
}

type FileRepository interface {
	Create(file *File) (*File, error)
	GetByID(id int) (*File, error)
}

type FileStorage interface {
	Save(filename string, src io.Reader) error
	Open(filename string) (io.ReadCloser, error)
	Delete(filename string) error
}
