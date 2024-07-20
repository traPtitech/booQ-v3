package router

import (
	"os"
	"testing"

	"github.com/traPtitech/booQ-v3/model"
)

func TestMain(m *testing.M) {
	model.SetupTestDB()
	exitCode := m.Run()
	os.Exit(exitCode)
}
