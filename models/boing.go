package models

import (
	"gorm.io/gorm"
	"time"
)

type Boing struct {
	gorm.Model
	Text      string    `json:"text"`
	UserId    uint      `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Boing) TableName() string {
	return "boings"
}

func NewBoing(text string, userId uint) Boing {
	return Boing{
		Text:      text,
		UserId:    userId,
		CreatedAt: time.Now(),
	}
}
