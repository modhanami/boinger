package services

import (
	"errors"
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var InvalidCredentials = "invalid credentials"

func setupDBForAuthService(t *testing.T) *gorm.DB {
	gdb := testutils.InitInMemDB(t)

	err := gdb.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		return nil
	}

	return gdb.Debug().Begin()
}

func TestAuthService_Authenticate(t *testing.T) {
	tests := []struct {
		name        string
		seed        func(db *gorm.DB)
		username    string
		password    string
		expectedErr error
	}{
		{
			name: "user found",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{Username: "user1", Password: "password1"})
			},
			username:    "user1",
			password:    "password1",
			expectedErr: nil,
		},
		{
			name:        "user not found",
			seed:        func(db *gorm.DB) {},
			username:    "user1",
			password:    "password1",
			expectedErr: ErrUserNotFound,
		},
		{
			name: "invalid credentials",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{Username: "user1", Password: "password1"})
			},
			username:    "user1",
			password:    "password2",
			expectedErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb := setupDBForAuthService(t)
			tt.seed(gdb)
			userService := NewUserService(gdb, logger.NewNoopLogger())
			tokenService := NewUserTokenService(gdb)
			service := NewAuthService(gdb, userService, tokenService, &fakePasswordHasher{})

			user, err := service.Authenticate(tt.username, tt.password)

			assert.ErrorIs(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.username, user.Username)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		seed        func(db *gorm.DB)
		username    string
		email       string
		password    string
		expectedErr error
	}{
		{
			name:        "user not found",
			seed:        func(db *gorm.DB) {},
			username:    "user1",
			email:       "user1@test.com",
			password:    "password1",
			expectedErr: nil,
		},
		{
			name: "user already exists",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{Username: "user1", Password: "password1"})
			},
			username:    "user1",
			email:       "user1@test.com",
			password:    "password1",
			expectedErr: ErrUserAlreadyExists,
		},
		{
			name:        "invalid credentials",
			seed:        func(db *gorm.DB) {},
			username:    "user1",
			email:       "user1@test.com",
			password:    InvalidCredentials,
			expectedErr: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb := setupDBForAuthService(t)
			tt.seed(gdb)
			userService := NewUserService(gdb, logger.NewNoopLogger())
			tokenService := NewUserTokenService(gdb)
			service := NewAuthService(gdb, userService, tokenService, &fakePasswordHasher{})

			user, err := service.Register(tt.username, tt.email, tt.password)

			assert.ErrorIs(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.username, user.Username)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

type fakePasswordHasher struct{}

func (n *fakePasswordHasher) HashPassword(password string) (string, error) {
	if password == InvalidCredentials {
		return "", errors.New("fakePasswordHasher: invalid credentials")
	}
	return password, nil
}

func (n *fakePasswordHasher) ComparePassword(hashedPassword, password string) error {
	if hashedPassword == password {
		return nil
	}

	return errors.New("fakePasswordHasher: password does not match")
}
