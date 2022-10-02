package models

type Boing struct {
	Model
	Text     string    `json:"text"`
	UserId   uint      `json:"-"`
	Comments []Comment `json:"comments"`
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
