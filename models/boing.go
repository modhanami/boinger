package models

import "time"

type BoingModel struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	Uid       string    `json:"id"`
	Text      string    `json:"text"`
	UserId    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *BoingModel) TableName() string {
	return "boings"
}

func NewBoing(uid string, text string, userId uint) BoingModel {
	return BoingModel{
		Uid:       uid,
		Text:      text,
		UserId:    userId,
		CreatedAt: time.Now(),
	}
}
