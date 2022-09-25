package models

type Boing struct {
	Base
	Text   string `json:"text"`
	UserId uint   `json:"-"`
}

func (m *Boing) TableName() string {
	return "boings"
}

func NewBoing(text string, userId uint) *Boing {
	return &Boing{
		Text:   text,
		UserId: userId,
	}
}
