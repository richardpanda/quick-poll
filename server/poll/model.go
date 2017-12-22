package poll

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/richardpanda/quick-poll/server/choice"
	"github.com/richardpanda/quick-poll/server/vote"
)

type Poll struct {
	ID        string          `json:"id" sql:"type:uuid;primary_key"`
	Question  string          `json:"question" gorm:"type:varchar(280)"`
	Choices   []choice.Choice `json:"choices"`
	CheckIP   bool            `json:"check_ip" gorm:"type:boolean" sql:"default:false"`
	Votes     []vote.Vote     `json:"votes"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *time.Time      `json:"-"`
}

type GetPollResponseBody struct {
	ID       string          `json:"id"`
	Question string          `json:"question"`
	Choices  []choice.Choice `json:"choices"`
	CheckIP  bool            `json:"check_ip"`
}

type PostPollsRequestBody struct {
	Question string   `json:"question"`
	Choices  []string `json:"choices"`
	CheckIP  bool     `json:"check_ip"`
}

type PostPollsResponseBody struct {
	ID       string          `json:"id"`
	Question string          `json:"question"`
	Choices  []choice.Choice `json:"choices"`
	CheckIP  bool            `json:"check_ip"`
}

func CreateTable(db *gorm.DB) {
	db.CreateTable(&Poll{})
}

func DropTable(db *gorm.DB) {
	db.DropTableIfExists("poll")
}
