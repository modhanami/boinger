package services

import (
	"github.com/modhanami/boinger/models"
	"github.com/modhanami/boinger/services/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestCommentService_Create(t *testing.T) {
	tests := []struct {
		name        string
		givenFunc   func(t *testing.T, db *gorm.DB)
		boingId     uint
		userId      uint
		text        string
		expectedErr error
	}{
		{
			name: "success",
			givenFunc: func(t *testing.T, db *gorm.DB) {
				user := models.User{Username: "test", Password: "test"}
				if err := db.Create(&user).Error; err != nil {
					t.Fatal(err)
				}

				boing := models.NewBoing("test boing", user.ID)
				if err := db.Create(&boing).Error; err != nil {
					t.Fatal(err)
				}
			},
			boingId:     1,
			userId:      1,
			text:        "test comment",
			expectedErr: nil,
		},
		{
			name: "boing not found",
			givenFunc: func(t *testing.T, db *gorm.DB) {
				user := models.User{Username: "test", Password: "test"}
				if err := db.Create(&user).Error; err != nil {
					t.Fatal(err)
				}
			},
			boingId:     1,
			userId:      1,
			text:        "test comment",
			expectedErr: ErrCommentCreationFailed,
		},
		{
			name: "user not found",
			givenFunc: func(t *testing.T, db *gorm.DB) {
				boing := models.NewBoing("test boing", 1)
				if err := db.Create(&boing).Error; err != nil {
					t.Fatal(err)
				}
			},
			boingId:     1,
			userId:      1,
			text:        "test comment",
			expectedErr: ErrCommentCreationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb := setupDBForCommentService(t)
			tt.givenFunc(t, gdb)

			s := NewCommentService(gdb)

			err := s.Create(tt.boingId, tt.userId, tt.text)
			assert.ErrorIs(t, tt.expectedErr, err)
		})
	}
}

func setupDBForCommentService(t *testing.T) *gorm.DB {
	gdb := testutils.InitInMemDB(t)

	err := gdb.AutoMigrate(&models.User{}, &models.Boing{}, &models.Comment{})
	if err != nil {
		return nil
	}

	return gdb.Debug().Begin()
}
