package storage

import (
	"io"

	"github.com/traPtitech/booQ-v3/domain"
)

// fileStorageAdapter は storage.Storage を domain.FileStorage に適合させるアダプター
type fileStorageAdapter struct {
	storage Storage
}

// NewFileStorage は domain.FileStorage を返します
func NewFileStorage() domain.FileStorage {
	return &fileStorageAdapter{storage: current}
}

func (a *fileStorageAdapter) Save(filename string, src io.Reader) error {
	return a.storage.Save(filename, src)
}

func (a *fileStorageAdapter) Open(filename string) (io.ReadCloser, error) {
	return a.storage.Open(filename)
}
