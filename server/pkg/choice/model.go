package choice

import (
	"time"

	"github.com/satori/go.uuid"
)

type Choice struct {
	ID       string `json:"id" sql:"type:uuid;primary_key"`
	PollID   string `json:"-" gorm:"type:uuid REFERENCES polls(id)"`
	Text     string `json:"text" gorm:"type:varchar(280)"`
	NumVotes int    `json:"num_votes"`
	// Vote      []vote.Vote `json:"votes"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type PostChoiceResponseBody struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	NumVotes int    `json:"num_votes"`
}

func New(text string) Choice {
	return Choice{
		ID:       uuid.NewV4().String(),
		Text:     text,
		NumVotes: 0,
	}
}
