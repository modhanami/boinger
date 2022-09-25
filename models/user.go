package models

type User struct {
	Base
	Username string `gorm:"not null;unique" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

func (m *User) TableName() string {
	return "users"
}
