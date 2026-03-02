package storage

import (
	"github.com/traPtitech/booQ-v3/domain"
)

var current domain.FileStorage = &Memory{files: map[string][]byte{}}

// NewFileStorage は domain.FileStorage を返します
func NewFileStorage() domain.FileStorage {
	return current
}
