package poll

import (
	"time"

	"github.com/richardpanda/quick-poll/server/pkg/choice"
)

type Poll struct {
	ID        string `sql:"type:uuid;primary_key"`
	Question  string `gorm:"type:varchar(280)"`
	Choices   []choice.Choice
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type POSTPollsRequestBody struct {
	Question string   `json:"question"`
	Choices  []string `json:"choices"`
}

type POSTPollsResponseBody struct {
	ID       string          `json:"id"`
	Question string          `json:"question"`
	Choices  []choice.Choice `json:"choices"`
}
