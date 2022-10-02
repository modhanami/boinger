package testutils

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/modhanami/boinger/logger"
	"gorm.io/driver/mysql"
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

func InitMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Error(err)
	}

	return db, mock
}

func initLogger() logger.Logger {
	return logger.NewNoopLogger()
}
