package services

import (
	"github.com/modhanami/boinger/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUserToken(t *testing.T) {
	var user = &models.UserModel{
		Id:       123,
		Uid:      "A1",
		Username: "bingbong",
		Password: "lookAtHimGo",
	}
	service := NewUserTokenService()

	token, _, err := service.Create(user, CreateOptions{})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyValidUserToken(t *testing.T) {
	var user = &models.UserModel{
		Id:       123,
		Uid:      "A1",
		Username: "bingbong",
		Password: "lookAtHimGo",
	}
	service := NewUserTokenService()
	token, _, err := service.Create(user, CreateOptions{})
	if err != nil {
		t.FailNow()
	}

	_, err = service.Verify(token)

	assert.NoError(t, err)
}

func TestFailVerifyExpiredUserToken(t *testing.T) {
	var user = &models.UserModel{
		Id:       123,
		Uid:      "A1",
		Username: "bingbong",
		Password: "lookAtHimGo",
	}
	service := NewUserTokenService()
	token, _, err := service.Create(user, CreateOptions{Exp: time.Now().Add(-time.Hour)})

	_, err = service.Verify(token)

	assert.Error(t, err)
}
