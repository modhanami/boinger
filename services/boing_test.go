package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoingService_List(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	var boings = []models.Boing{
		models.NewBoing("boing1", 0),
		models.NewBoing("boing2", 0),
		models.NewBoing("boing3", 0),
	}

	rows := createBoingRows()

	for _, boing := range boings {
		rows = rows.AddRow(boing.Id, boing.Text, boing.UserId, boing.CreatedAt)
	}

	mock.ExpectQuery("SELECT *").WillReturnRows(rows)

	boingsResult, err := service.List()

	assert.NoError(t, err)
	assert.Equal(t, boingsResult, boings)
}

func TestBoingService_List_UnexpectedDBError(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	mock.ExpectQuery("SELECT").WillReturnError(ErrUnexpectedDBError)

	_, err := service.List()

	assert.Error(t, err)
	assert.Equal(t, err, ErrUnexpectedDBError)
}

func TestBoingService_Create(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	boing := models.NewBoing("boing1", 0)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := service.Create(boing.Text, 0)

	assert.NoError(t, err)
}

func TestBoingService_Create_Failed(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WillReturnError(ErrBoingCreationFailed)
	mock.ExpectRollback()

	err := service.Create("", 0)

	assert.Error(t, err)
	assert.Equal(t, err, ErrBoingCreationFailed)
}

func TestBoingService_GetById(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	boing := models.NewBoing("boing1", 0)

	rows := createBoingRows()
	rows.AddRow(boing.Id, boing.Text, boing.UserId, boing.CreatedAt)

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	boingResult, err := service.GetById(boing.Id)

	assert.NoError(t, err)
	assert.Equal(t, boingResult, boing)
}

func TestBoingService_GetById_NotFound(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	mock.ExpectQuery("SELECT").WillReturnRows(createBoingRows())

	_, err := service.GetById(1)

	assert.Error(t, err)
	assert.Equal(t, err, ErrBoingNotFound)
}

func TestBoingService_GetById_UnexpectedDBError(t *testing.T) {
	service, mock := initBoingServiceWithMocks(t)

	mock.ExpectQuery("SELECT").WillReturnError(ErrUnexpectedDBError)

	_, err := service.GetById(1)

	assert.Error(t, err)
	assert.Equal(t, err, ErrUnexpectedDBError)
}

func createBoingRows() *sqlmock.Rows {
	var columns = []string{"uid", "text", "user_id", "created_at"}
	var rows = sqlmock.NewRows(columns)
	return rows
}

func initBoingServiceWithMocks(t *testing.T) (BoingService, sqlmock.Sqlmock) {
	db, mock := initMockDB(t)
	return NewBoingService(db, logger.NewNoopLogger()), mock
}
