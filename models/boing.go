package models

import "time"

type Boing struct {
	Id        uint      `gorm:"primary_key" json:"-"`
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
