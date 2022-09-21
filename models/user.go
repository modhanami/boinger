package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func (m *User) TableName() string {
	return "users"
}
