package services

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	var (
		sqlDB *sql.DB
		mock  sqlmock.Sqlmock
		err   error
	)

	sqlDB, mock, err = sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	return db, mock, nil
}
