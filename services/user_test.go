package services

import (
	"github.com/modhanami/boinger/logger"
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func setupDBForUserService(t *testing.T) *gorm.DB {
	gdb, err := testutils.InitInMemDB(t)

	err = gdb.AutoMigrate(&models.User{})
	if err != nil {
		return nil
	}

	return gdb.Debug().Begin()
}

func TestUserService_Create(t *testing.T) {
	tests := []struct {
		name        string
		seed        func(db *gorm.DB)
		user        *models.User
		expectedErr error
	}{
		{
			name: "create user",
			seed: func(db *gorm.DB) {},
			user: &models.User{
				Username: "user1",
				Email:    "email1@test.com",
				Password: "password1",
			},
			expectedErr: nil,
		},
		{
			name: "create user with duplicate username",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{
					Username: "user1",
					Email:    "email1@test.com",
					Password: "password1",
				})
			},
			user: &models.User{
				Username: "user1",
				Email:    "email2@test.com",
				Password: "password1",
			},
			expectedErr: ErrUserCreationFailed,
		},
		{
			name: "create user with duplicate email",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{
					Username: "user1",
					Email:    "email1@test.com",
					Password: "password1",
				})
			},
			user: &models.User{
				Username: "user2",
				Email:    "email1@test.com",
				Password: "password1",
			},
			expectedErr: ErrUserCreationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDBForUserService(t)
			tt.seed(db)
			service := NewUserService(db, logger.NewNoopLogger())

			user, err := service.Create(tt.user)

			assert.ErrorIs(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				assert.NotEmpty(t, user.ID)
			}
		})
	}
}

func TestUserService_ExistsByUsername(t *testing.T) {
	tests := []struct {
		name        string
		seed        func(db *gorm.DB)
		username    string
		expectedErr error
		expected    bool
	}{
		{
			name: "exists",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{
					Username: "user1",
					Email:    "email1@test.com",
					Password: "password1",
				})
			},
			username:    "user1",
			expectedErr: nil,
			expected:    true,
		},
		{
			name:        "does not exist",
			seed:        func(db *gorm.DB) {},
			username:    "user1",
			expectedErr: nil,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDBForUserService(t)
			tt.seed(db)
			service := NewUserService(db, logger.NewNoopLogger())

			exists, err := service.ExistsByUsername(tt.username)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestUserService_GetById(t *testing.T) {
	tests := []struct {
		name        string
		seed        func(db *gorm.DB)
		id          uint
		expectedErr error
		expected    *models.User
	}{
		{
			name: "found",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{
					Username: "user1",
					Email:    "email1@test.com",
					Password: "password1",
				})
			},
			id:          1,
			expectedErr: nil,
		},
		{
			name:        "not found",
			seed:        func(db *gorm.DB) {},
			id:          1,
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDBForUserService(t)
			tt.seed(db)
			service := NewUserService(db, logger.NewNoopLogger())

			user, err := service.GetById(tt.id)

			assert.ErrorIs(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.id, user.ID)
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	tests := []struct {
		name        string
		seed        func(db *gorm.DB)
		username    string
		expectedErr error
		expected    *models.User
	}{
		{
			name: "found",
			seed: func(db *gorm.DB) {
				db.Create(&models.User{
					Username: "user1",
					Email:    "email1@test.com",
					Password: "password1",
				})
			},
			username:    "user1",
			expectedErr: nil,
		},
		{
			name:        "not found",
			seed:        func(db *gorm.DB) {},
			username:    "user1",
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDBForUserService(t)
			tt.seed(db)
			service := NewUserService(db, logger.NewNoopLogger())

			user, err := service.GetByUsername(tt.username)

			assert.ErrorIs(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				assert.Equal(t, tt.username, user.Username)
			}
		})
	}
}
