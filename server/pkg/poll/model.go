package poll

import (
	"time"

	"github.com/richardpanda/quick-poll/server/pkg/choice"
)

type Poll struct {
	ID        string          `json:"id" sql:"type:uuid;primary_key"`
	Question  string          `json:"question" gorm:"type:varchar(280)"`
	Choices   []choice.Choice `json:"choices"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *time.Time      `json:"-"`
}

type GETPollResponseBody struct {
	ID       string          `json:"id"`
	Question string          `json:"question"`
	Choices  []choice.Choice `json:"choices"`
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
