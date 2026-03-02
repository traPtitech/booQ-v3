package storage

import (
	"bytes"
	"io"
	"sync"

	"github.com/traPtitech/booQ-v3/domain"
)

// Memory インメモリストレージ
type Memory struct {
	mu    sync.RWMutex
	files map[string][]byte
}

// SetMemoryStorage メモリストレージをカレントストレージに設定します
func SetMemoryStorage() {
	current = &Memory{files: map[string][]byte{}}
}

func (m *Memory) Save(filename string, src io.Reader) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	m.files[filename] = b
	return nil
}

func (m *Memory) Open(filename string) (io.ReadCloser, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.files[filename]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return io.NopCloser(bytes.NewReader(b)), nil
}

func (m *Memory) Delete(filename string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.files[filename]
	if !ok {
		return domain.ErrNotFound
	}
	delete(m.files, filename)
	return nil
}
