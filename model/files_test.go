package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTableName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "files", (&File{}).TableName())
}

func TestCreateFile(t *testing.T) {
	t.Parallel()
	id := "testUser"

	f, err := CreateFile(id, strings.NewReader("test file"), "txt")
	if assert.NoError(t, err) {
		assert.NotEmpty(t, f.ID)
		assert.Equal(t, id, f.UploadUserID)
	}
}
