package models

import "time"

type RefreshToken struct {
	Base
	UserId    uint `gorm:"not null"`
	User      User
	Token     string `gorm:"not null"`
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
