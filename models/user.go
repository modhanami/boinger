package models

import "time"

type User struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	Uid       string    `json:"uid"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *User) TableName() string {
	return "users"
}

func NewUser(Uid string, Username string, Password string) User {
	return User{
		Uid:       Uid,
		Username:  Username,
		Password:  Password,
		CreatedAt: time.Now(),
	}
}
