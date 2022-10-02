package models

type Comment struct {
	Model
	BoingID uint   `gorm:"not null" json:"-"`
	Boing   Boing  `json:"-"`
	UserID  uint   `gorm:"not null" json:"-"`
	User    User   `json:"-"`
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
