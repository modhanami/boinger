package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/modhanami/boinger/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func initMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func initLogger() log.Interface {
	return log.NewNoop()
}
