package testutils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func InitInMemDB(t *testing.T) (*gorm.DB, error) {
	gdb, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fail()
	}
	return gdb, err
}
