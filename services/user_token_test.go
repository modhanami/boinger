package services

import (
	"github.com/modhanami/boinger/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUserToken(t *testing.T) {
	gdb := setup(t)
	var user = &models.User{
		Id:       123,
		Uid:      "A1",
		Username: "bingbong",
		Password: "lookAtHimGo",
	}
	service := NewUserTokenService(gdb)

	token, _, err := service.Create(user, CreateOptions{})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyValidUserToken(t *testing.T) {
	gdb := setup(t)
	var user = &models.User{
		Id:       123,
		Uid:      "A1",
		Username: "bingbong",
		Password: "lookAtHimGo",
	}
	service := NewUserTokenService(gdb)
	token, _, err := service.Create(user, CreateOptions{})
	if err != nil {
		t.FailNow()
	}

	userClaims, err := service.Verify(token)

	assert.NoError(t, err)
	assert.NotEmpty(t, userClaims)
	assert.Equal(t, user.Uid, userClaims.Uid)
}

func TestFailVerifyExpiredUserToken(t *testing.T) {
	gdb := setup(t)
	var user = &models.User{
		Id:       123,
		Uid:      "A1",
		Username: "bingbong",
		Password: "lookAtHimGo",
	}
	service := NewUserTokenService(gdb)
	token, _, err := service.Create(user, CreateOptions{Exp: time.Now().Add(-time.Hour)})

	userClaims, err := service.Verify(token)

	assert.Error(t, err)
	assert.Empty(t, userClaims)
}

func TestRenewRefreshToken(t *testing.T) {
	gdb := setup(t)
	service := NewUserTokenService(gdb)
	user := models.User{Uid: "uid1", Username: "user1", Password: "password1"}
	gdb.Create(&user)

	token, err := service.RenewRefreshToken(user.Id)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token.Token, 64)
	assert.Equal(t, user.Id, token.UserId)
}

func TestRenewRefreshToken_RevokesOldRefreshToken(t *testing.T) {
	gdb := setup(t)
	service := NewUserTokenService(gdb)
	user := models.User{Uid: "uid1", Username: "user1", Password: "password1"}
	gdb.Create(&user)
	oldToken, err := service.RenewRefreshToken(user.Id)
	if err != nil {
		t.FailNow()
	}

	newToken, err := service.RenewRefreshToken(user.Id)

	assert.NoError(t, err)
	assert.NotEqual(t, oldToken.Token, newToken.Token)

	var count int64
	gdb.Model(&models.RefreshToken{}).Where("revoked_at IS NOT NULL").Count(&count)
	assert.Equal(t, int64(1), count)
}
