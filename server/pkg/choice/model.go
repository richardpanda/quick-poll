package choice

import (
	"time"

	"github.com/satori/go.uuid"
)

type Choice struct {
	ID        string     `sql:"type:uuid;primary_key" json:"id"`
	PollID    string     `gorm:"type:uuid REFERENCES polls(id)" json:"-"`
	Text      string     `gorm:"type:varchar(280)" json:"text"`
	NumVotes  int        `json:"num_votes"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func New(text string) Choice {
	return Choice{
		ID:       uuid.NewV4().String(),
		Text:     text,
		NumVotes: 0,
	}
}
