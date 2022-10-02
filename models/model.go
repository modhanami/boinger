package models

import "time"

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
