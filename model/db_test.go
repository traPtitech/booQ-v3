package model

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetUpTestDB()
	exitCode := m.Run()
	os.Exit(exitCode)
}
