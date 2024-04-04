package model

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetupTestDB()
	exitCode := m.Run()
	os.Exit(exitCode)
}
