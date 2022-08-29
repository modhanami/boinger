package services

import (
	"errors"
	"fmt"
	"github.com/modhanami/boinger/log"
	"github.com/modhanami/boinger/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

var InvalidCredentials = "invalid credentials"

func setup(t *testing.T) *gorm.DB {
	gdb, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fail()
	}

	err = gdb.AutoMigrate(&models.User{})
	if err != nil {
		return nil
	}

	return gdb.Debug().Begin()
}

func TestAuthService_Authenticate(t *testing.T) {
	gdb := setup(t)
	gdb.Create(&models.User{Uid: "uid1", Username: "user1", Password: "password1"})
	service := NewAuthService(gdb, NewUserService(gdb, log.NewNoop()), NewUserTokenService(), &fakePasswordHasher{})

	user, err := service.Authenticate("user1", "password1")

	assert.NoError(t, err)
	assert.Equal(t, "uid1", user.Uid)
}

func TestAuthService_Authenticate_UserNotFound(t *testing.T) {
	gdb := setup(t)
	service := NewAuthService(gdb, NewUserService(gdb, log.NewNoop()), NewUserTokenService(), &fakePasswordHasher{})

	user, err := service.Authenticate("user1", "password1")

	assert.Equal(t, err, ErrUserNotFound)
	assert.Empty(t, user)

}

func TestAuthService_Authenticate_InvalidCredentials(t *testing.T) {
	gdb := setup(t)
	gdb.Create(&models.User{Uid: "uid1", Username: "user1", Password: "password1"})
	service := NewAuthService(gdb, NewUserService(gdb, log.NewNoop()), NewUserTokenService(), &fakePasswordHasher{})

	user, err := service.Authenticate("user1", "password2")

	fmt.Println(err)
	assert.Equal(t, err, ErrInvalidCredentials)
	assert.Empty(t, user)
}

func TestAuthService_Register(t *testing.T) {
	gdb := setup(t)
	service := NewAuthService(gdb, NewUserService(gdb, log.NewNoop()), NewUserTokenService(), &fakePasswordHasher{})

	user, err := service.Register("user1", "password1")

	assert.NoError(t, err)
	assert.Equal(t, "user1", user.Username)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	gdb := setup(t)
	gdb.Create(&models.User{Uid: "uid1", Username: "user1", Password: "password1"})
	service := NewAuthService(gdb, NewUserService(gdb, log.NewNoop()), NewUserTokenService(), &fakePasswordHasher{})

	user, err := service.Register("user1", "password1")

	assert.Equal(t, err, ErrUserAlreadyExists)
	assert.Empty(t, user)
}

func TestAuthService_Register_InvalidCredentials(t *testing.T) {
	gdb := setup(t)
	service := NewAuthService(gdb, NewUserService(gdb, log.NewNoop()), NewUserTokenService(), &fakePasswordHasher{})

	user, err := service.Register("user1", InvalidCredentials)

	assert.Equal(t, err, ErrInvalidCredentials)
	assert.Empty(t, user)
}

type fakePasswordHasher struct{}

func (n *fakePasswordHasher) HashPassword(password string) (string, error) {
	if password == InvalidCredentials {
		return "", ErrInvalidCredentials
	}
	return password, nil
}

func (n *fakePasswordHasher) ComparePassword(hashedPassword, password string) error {
	if hashedPassword == password {
		return nil
	}

	return errors.New("fakePasswordHasher: password does not match")
}

//func initAuthServiceWithSuccessMocks(t *testing.T) (AuthService, sqlmock.Sqlmock) {
//	db, mock := initMockDB(t)
//	return NewAuthService(db, &mockUserServiceSuccess{}, &mockUserTokenService{}, &mockPasswordHasher{}), mock
//}
//
//func initAuthServiceWithErrorMocks(t *testing.T) (AuthService, sqlmock.Sqlmock) {
//	db, mock := initMockDB(t)
//	return NewAuthService(db, &mockUserServiceError{}, &mockUserTokenService{}, &mockPasswordHasher{}), mock
//}
//
//type mockPasswordHasher struct{}
//
//func (m *mockPasswordHasher) HashPassword(password string) (string, error) {
//	return password, nil
//}
//
//func (m *mockPasswordHasher) ComparePassword(string, string) error {
//	return nil
//}
//
//type mockUserTokenService struct {
//	UserTokenService
//}
//
//type mockUserServiceSuccess struct {
//	UserService
//}
//
//func (m *mockUserServiceSuccess) GetByUsername(username string) (models.User, error) {
//	return models.User{
//		Uid:       "uid",
//		Username:  username,
//		Password:  "password",
//		CreatedAt: time.Date(2020, time.May, 5, 8, 0, 0, 0, time.UTC),
//	}, nil
//}
//
//func (m *mockUserServiceSuccess) ExistsByUsername(string) (bool, error) {
//	return false, nil
//}
//
//func (m *mockUserServiceSuccess) Create(username string, password string) (models.User, error) {
//	return models.User{
//		Uid:       "uid",
//		Username:  username,
//		Password:  password,
//		CreatedAt: time.Date(2020, time.May, 5, 8, 0, 0, 0, time.UTC),
//	}, nil
//}
//
//type mockUserServiceError struct {
//	UserService
//}
//
//func (m *mockUserServiceError) GetByUsername(username string) (models.User, error) {
//	return models.User{}, ErrUserNotFound
//}
