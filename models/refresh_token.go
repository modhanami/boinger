package models

import "time"

type RefreshToken struct {
	Id        uint `gorm:"primary_key"`
	UserId    uint
	User      User
	Token     string
	RevokedAt *time.Time
}

func (m *RefreshToken) TableName() string {
	return "refresh_tokens"
}

func NewRefreshToken(userId uint, token string) *RefreshToken {
	return &RefreshToken{
		UserId: userId,
		Token:  token,
	}
}
