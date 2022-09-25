package models

type Comment struct {
	Base
	BoingID uint `gorm:"not null" json:"boingId"`
	Boing   Boing
	UserID  uint `gorm:"not null" json:"userId"`
	User    User
	Text    string `json:"text"`
}

func (m *Comment) TableName() string {
	return "comments"
}

func NewComment(boingId, userId uint, text string) *Comment {
	return &Comment{
		BoingID: boingId,
		UserID:  userId,
		Text:    text,
	}
}
