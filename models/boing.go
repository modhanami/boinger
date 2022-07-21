package models

import "time"

type Boing struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	Uid       string    `json:"id"`
	Text      string    `json:"text"`
	UserId    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Boing) TableName() string {
	return "boings"
}

func NewBoing(uid string, text string, userId uint) Boing {
	return Boing{
		Uid:       uid,
		Text:      text,
		UserId:    userId,
		CreatedAt: time.Now(),
	}
}
