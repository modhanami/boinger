package services

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/modhanami/boinger/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestShouldListBoings(t *testing.T) {
	db, mock, err := initMockDB()
	if err != nil {
		t.Error(err)
	}

	var boings = []models.BoingModel{
		NewBoing("A1", "boing1", 0),
		NewBoing("A2", "boing2", 0),
		NewBoing("A3", "boing3", 0),
	}

	rows := createRows()

	for _, boing := range boings {
		rows = rows.AddRow(boing.Uid, boing.Text, boing.UserId, boing.CreatedAt)
	}

	mock.ExpectQuery("SELECT *").WillReturnRows(rows)

	var service = NewBoingService(db)
	var boingsResult []models.BoingModel
	boingsResult, err = service.List()
	if err != nil {
		t.Error(err)
	}

	if len(boingsResult) != len(boings) {
		t.Fatalf("expected %d boings, got %d", len(boings), len(boingsResult))
	}

	for i, boing := range boingsResult {
		if !reflect.DeepEqual(boing, boings[i]) {
			t.Fatalf("expected %v, got %v", boings[i], boing)
		}
	}
}

func TestShouldCreateBoing(t *testing.T) {
	db, mock, err := initMockDB()
	if err != nil {
		t.Error(err)
	}

	boing := NewBoing("A1", "boing1", 0)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	service := NewBoingService(db)
	err = service.Create(boing.Text, 0)
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestShouldGetBoingById(t *testing.T) {
	db, mock, err := initMockDB()
	if err != nil {
		t.Error(err)
	}

	boing := NewBoing("A1", "boing1", 0)

	rows := createRows()
	rows.AddRow(boing.Uid, boing.Text, boing.UserId, boing.CreatedAt)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	var service = NewBoingService(db)
	var boingResult models.BoingModel
	boingResult, err = service.Get(boing.Id)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(boingResult, boing) {
		t.Fatalf("expected %v, got %v", boing, boingResult)
	}
}

func createRows() *sqlmock.Rows {
	var columns = []string{"uid", "text", "user_id", "created_at"}
	var rows = sqlmock.NewRows(columns)
	return rows
}

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
