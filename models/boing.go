package models

import "time"

type Boing struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	Uid       string    `json:"uid"`
	Text      string    `json:"text"`
	UserId    uint      `json:"-"`
	UserUid   string    `json:"userUid"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Boing) TableName() string {
	return "boings"
}

func NewBoing(uid string, text string, userId uint, userUid string) Boing {
	return Boing{
		Uid:       uid,
		Text:      text,
		UserId:    userId,
		UserUid:   userUid,
		CreatedAt: time.Now(),
	}
}
