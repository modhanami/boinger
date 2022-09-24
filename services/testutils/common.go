package testutils

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func InitInMemDB(t *testing.T) *gorm.DB {
	gdb, err := gorm.Open(sqlite.Open("file::memory:?_foreign_keys=true"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return gdb
}
