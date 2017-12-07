package choice

import (
	"time"
)

type Choice struct {
	ID        string     `sql:"type:uuid;primary_key" json:"id"`
	PollID    string     `gorm:"type:uuid REFERENCES polls(id)" json:"-"`
	Text      string     `gorm:"type:varchar(280)" json:"text"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
