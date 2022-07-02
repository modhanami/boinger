package models

import "time"

type UserModel struct {
	Id        uint `gorm:"primary_key" json:"-"`
	Uid       string
	Username  string
	Password  string
	CreatedAt time.Time
}

func (m *UserModel) TableName() string {
	return "users"
}
