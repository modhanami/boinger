package services

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/modhanami/boinger/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestUserService_Create(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	user := models.User{
		Username: "bingbong",
		Password: "eeur",
	}

	rows := createUserRows()
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user, err := service.Create(user.Username, user.Password)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, "bingbong")
	assert.NotEqual(t, user.Password, "eeur")
}

func TestUserService_Create_DuplicateUsername(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	user := models.User{
		Username: "bingbong",
		Password: "eeur",
	}

	rows := createUserRows()
	rows.AddRow(user.ID, user.Username, user.Password, user.CreatedAt)
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	_, err := service.Create(user.Username, user.Password)

	assert.Error(t, err)
	assert.Equal(t, err, ErrUserAlreadyExists)
}

func TestUserService_Exists_Found(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	user := models.User{
		Username: "bingbong",
		Password: "eeur",
	}

	rows := createUserRows()
	rows.AddRow(user.ID, user.Username, user.Password, user.CreatedAt)
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	exists, err := service.ExistsByUsername(user.Username)

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserService_Exists_NotFound(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	rows := createUserRows()
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	exists, err := service.ExistsByUsername("whosthis")

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserService_GetById_Found(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	user := models.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "bingbong",
		Password: "eeur",
	}

	rows := createUserRowsWithId()
	rows.AddRow(1, user.ID, user.Username, user.Password, user.CreatedAt)
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	user, err := service.GetById(1)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, "A1")
}

func TestUserService_GetById_NotFound(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	rows := createUserRows()
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	_, err := service.GetById(1)

	assert.Error(t, err)
	assert.Equal(t, err, ErrUserNotFound)
}

func TestUserService_GetByUsername_Found(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	user := models.User{
		Username: "bingbong",
		Password: "eeur",
	}

	rows := createUserRows()
	rows.AddRow(user.ID, user.Username, user.Password, user.CreatedAt)
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	user, err := service.GetByUsername(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, "A1")
}

func TestUserService_GetByUsername_NotFound(t *testing.T) {
	service, mock := initServiceWithMocks(t)

	rows := createUserRows()
	mock.ExpectQuery("SELECT ").WillReturnRows(rows)

	_, err := service.GetByUsername("whosthis")

	assert.Error(t, err)
	assert.Equal(t, err, ErrUserNotFound)
}

func createUserRows() *sqlmock.Rows {
	var columns = []string{"uid", "username", "password", "created_at"}
	var rows = sqlmock.NewRows(columns)
	return rows
}

func createUserRowsWithId() *sqlmock.Rows {
	var columns = []string{"id", "uid", "username", "password", "created_at"}
	var rows = sqlmock.NewRows(columns)
	return rows
}

func initServiceWithMocks(t *testing.T) (UserService, sqlmock.Sqlmock) {
	db, mock := initMockDB(t)
	serviceLogger := initLogger()
	service := NewUserService(db, serviceLogger)

	return service, mock
}
